import datetime
import logging
import json

from airflow.models.dag import DAG
from airflow.operators.python import PythonOperator
from airflow.providers.google.cloud.hooks.gcs import GCSHook
from airflow.hooks.base import BaseHook
from airflow.providers.google.cloud.transfers.gcs_to_bigquery import GCSToBigQueryOperator

try:
    from kafka import KafkaConsumer
except ImportError:
    logging.error("Thư viện 'kafka-python' chưa được cài đặt. Hãy thêm nó vào môi trường Airflow.")
    raise

# =============================================================================
# CẤU HÌNH
# =============================================================================
KAFKA_CONN_ID = "kafka_redpanda_cloud"
KAFKA_TOPIC = "trip_search"

GCP_CONN_ID = "google_cloud_default"
GCS_BUCKET_NAME = "asia-southeast1-invoice-dua-e3630dcb-bucket"

BIGQUERY_PROJECT_ID = "dacntt-dfabb"
BIGQUERY_DATASET_NAME = "duancntt"
BIGQUERY_TABLE_NAME = "trip_searches_raw"

# =============================================================================
# HÀM PYTHON CHÍNH
# =============================================================================
def consume_kafka_to_gcs(**context):
    try:
        kafka_conn = BaseHook.get_connection(KAFKA_CONN_ID)
        bootstrap_servers = kafka_conn.host.split(',')
        username = kafka_conn.login
        password = kafka_conn.password
        
        logging.info(f"Đang kết nối đến Kafka topic: '{KAFKA_TOPIC}'...")

        consumer = KafkaConsumer(
            KAFKA_TOPIC,
            bootstrap_servers=bootstrap_servers,
            security_protocol="SASL_SSL",
            sasl_mechanism="SCRAM-SHA-256",
            sasl_plain_username=username,
            sasl_plain_password=password,
            group_id="airflow-trip-search-consumer-group",
            auto_offset_reset="earliest",
            consumer_timeout_ms=60000,
            request_timeout_ms=40000,
            session_timeout_ms=30000,
            heartbeat_interval_ms=10000
        )

        messages = []
        for message in consumer:
            try:
                msg_data = json.loads(message.value.decode('utf-8'))
                messages.append(msg_data)
            except json.JSONDecodeError:
                logging.warning(f"Bỏ qua message không phải JSON: {message.value}")

        consumer.close()
        logging.info(f"Đã đóng kết nối Kafka.")

        if not messages:
            logging.info("Không có message mới nào trong batch này. Kết thúc task.")
            context['ti'].xcom_push(key='gcs_object_path', value=None)
            return

        logging.info(f"Đã nhận được {len(messages)} messages.")
        
        ndjson_data = "\n".join([json.dumps(msg) for msg in messages])

        gcs_hook = GCSHook(gcp_conn_id=GCP_CONN_ID)
        execution_ts_nodash = context["ts_nodash"]
        object_path = f"data/kafka/{KAFKA_TOPIC}/{KAFKA_TOPIC}_{execution_ts_nodash}.json"
        
        gcs_hook.upload(
            bucket_name=GCS_BUCKET_NAME,
            object_name=object_path,
            data=ndjson_data.encode("utf-8"),
            mime_type="application/json"
        )
        logging.info(f"✅ Đã tải thành công batch lên GCS: gs://{GCS_BUCKET_NAME}/{object_path}")
        
        context['ti'].xcom_push(key='gcs_object_path', value=object_path)

    except Exception as e:
        logging.error(f"Lỗi trong quá trình đọc Kafka: {e}")
        raise

# =============================================================================
# SCHEMA CỦA BẢNG BIGQUERY (ĐÃ SỬA LẠI CHO ĐÚNG)
# =============================================================================
TRIP_SEARCH_SCHEMA = [
    {"name": "fromProvinceId", "type": "INTEGER", "mode": "NULLABLE"},
    {"name": "toProvinceId", "type": "INTEGER", "mode": "NULLABLE"},
    {"name": "departureDate", "type": "DATE", "mode": "NULLABLE"},
    {"name": "searchTimestamp", "type": "TIMESTAMP", "mode": "NULLABLE"},
    {"name": "quantity", "type": "INTEGER", "mode": "NULLABLE"},
    {"name": "userId", "type": "STRING", "mode": "NULLABLE"},
]

# =============================================================================
# ĐỊNH NGHĨA DAG
# =============================================================================
with DAG(
    dag_id='kafka_to_bigquery_pipeline_final',
    start_date=datetime.datetime(2025, 6, 26),
    schedule_interval='0 */2 * * *',
    catchup=False,
    tags=['kafka', 'gcs', 'bigquery', 'production'],
    default_args={
        'owner': 'data-team',
        'retries': 1,
        'retry_delay': datetime.timedelta(minutes=5),
    }
) as dag:

    consume_from_kafka_task = PythonOperator(
        task_id='consume_from_kafka_to_gcs',
        python_callable=consume_kafka_to_gcs,
    )

    load_gcs_to_bq_task = GCSToBigQueryOperator(
        task_id='load_gcs_to_bigquery',
        bucket=GCS_BUCKET_NAME,
        source_objects=["{{ ti.xcom_pull(task_ids='consume_from_kafka_to_gcs', key='gcs_object_path') }}"],
        trigger_rule='all_success',
        destination_project_dataset_table=f"{BIGQUERY_PROJECT_ID}.{BIGQUERY_DATASET_NAME}.{BIGQUERY_TABLE_NAME}",
        schema_fields=TRIP_SEARCH_SCHEMA,
        source_format='NEWLINE_DELIMITED_JSON',
        write_disposition='WRITE_APPEND',
        create_disposition='CREATE_IF_NEEDED',
        gcp_conn_id=GCP_CONN_ID,
    )

    consume_from_kafka_task >> load_gcs_to_bq_task