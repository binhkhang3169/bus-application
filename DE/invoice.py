"""
DAG ETL hàng ngày từ PostgreSQL sang GCS, và tải vào BigQuery.

VERSION 2 - Sửa lỗi:
1. Sửa lỗi chữ hoa/thường cho bảng `ticket` của PostgreSQL.
2. Sửa lỗi sai định dạng `destination_project_dataset_table` cho BigQuery.
"""
import datetime
import logging
import pandas as pd

from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.providers.google.cloud.hooks.gcs import GCSHook
from airflow.providers.postgres.hooks.postgres import PostgresHook
from airflow.providers.google.cloud.transfers.gcs_to_bigquery import GCSToBigQueryOperator

# =============================================================================
# Cấu hình chung
# =============================================================================
GCP_CONN_ID = "google_cloud_default"
GCS_BUCKET_NAME = "asia-southeast1-invoice-dua-e3630dcb-bucket"
POSTGRES_TICKET_CONN_ID = "postgres_ticket_db"
POSTGRES_INVOICE_CONN_ID = "postgres_invoice_db"

# Cấu hình cho BigQuery
# SỬA LỖI 2: Đảm bảo Project ID chỉ có một phần, không bị lặp.
BIGQUERY_PROJECT_ID = "dacntt-dfabb"       # <-- SỬA Ở ĐÂY
BIGQUERY_DATASET_NAME = "duancntt"             # <-- THAY THẾ bằng Dataset đúng của bạn

def extract_to_gcs(
    postgres_conn_id: str,
    table_name: str,
    time_column: str,
    gcs_bucket_name: str,
    schema: str = "public",
    **context
):
    """Hàm trích xuất dữ liệu từ Postgres và lưu vào GCS."""
    start_date = context["data_interval_start"]
    end_date = context["data_interval_end"]
    execution_date_str = context["ds_nodash"]

    object_path = f"data/{table_name.lower()}/{table_name.lower()}_{execution_date_str}.csv"
    context['ti'].xcom_push(key='gcs_object_path', value=object_path)

    logging.info(f"Bắt đầu xử lý bảng '{table_name}' cho khoảng: {start_date} -> {end_date}")

    # SỬA LỖI 1: Bỏ dấu ngoặc kép "" bao quanh tên bảng và cột để PostgreSQL tự xử lý.
    # Điều này giúp tương thích với cả tên viết hoa và viết thường.
    sql = (
        f'SELECT * FROM {schema}.{table_name} '
        f"WHERE {time_column} >= '{start_date.isoformat()}' "
        f"AND {time_column} < '{end_date.isoformat()}'"
    )
    logging.info(f"Executing SQL: {sql}")

    pg_hook = PostgresHook(postgres_conn_id=postgres_conn_id)
    df = pg_hook.get_pandas_df(sql)

    if df.empty:
        logging.info(f"Không có dữ liệu mới trong bảng '{table_name}'.")
        context['ti'].xcom_push(key='gcs_object_path', value=None)
        return

    logging.info(f"Đã đọc {len(df)} dòng từ bảng '{table_name}'.")

    gcs_hook = GCSHook(gcp_conn_id=GCP_CONN_ID)
    csv_data = df.to_csv(index=False, encoding="utf-8")

    gcs_hook.upload(
        bucket_name=gcs_bucket_name,
        object_name=object_path,
        data=csv_data.encode("utf-8"),
        mime_type="text/csv",
    )
    logging.info(f"✅ Tải thành công file lên GCS: gs://{gcs_bucket_name}/{object_path}")


with DAG(
    dag_id='etl_postgres_gcs_bigquery_daily_v2', # Đổi tên dag_id để tránh nhầm lẫn
    start_date=datetime.datetime(2025, 6, 26),
    schedule_interval='0 20 * * *',
    catchup=False,
    tags=['etl', 'gcs', 'bigquery', 'production'],
    default_args={'owner': 'data-team', 'retries': 1}
) as dag:

    # --- LUỒNG XỬ LÝ CHO BẢNG TICKET ---
    extract_tickets_to_gcs_task = PythonOperator(
        task_id='extract_tickets_to_gcs',
        python_callable=extract_to_gcs,
        op_kwargs={
            "postgres_conn_id": POSTGRES_TICKET_CONN_ID,
            # SỬA LỖI 1: Dùng tên chữ thường như trong database gợi ý.
            "table_name": "ticket",
            "time_column": "created_at",
            "gcs_bucket_name": GCS_BUCKET_NAME,
        },
    )

    load_tickets_to_bq_task = GCSToBigQueryOperator(
        task_id='load_tickets_to_bq',
        bucket=GCS_BUCKET_NAME,
        source_objects=["{{ ti.xcom_pull(task_ids='extract_tickets_to_gcs', key='gcs_object_path') }}"],
        destination_project_dataset_table=f"{BIGQUERY_PROJECT_ID}.{BIGQUERY_DATASET_NAME}.tickets",
        source_format='CSV',
        skip_leading_rows=1,
        write_disposition='WRITE_APPEND',
        autodetect=True,
        gcp_conn_id=GCP_CONN_ID,
    )

    # --- LUỒNG XỬ LÝ CHO BẢNG INVOICES ---
    extract_invoices_to_gcs_task = PythonOperator(
        task_id='extract_invoices_to_gcs',
        python_callable=extract_to_gcs,
        op_kwargs={
            "postgres_conn_id": POSTGRES_INVOICE_CONN_ID,
            "table_name": "invoices",
            "time_column": "created_at",
            "gcs_bucket_name": GCS_BUCKET_NAME,
        },
    )

    load_invoices_to_bq_task = GCSToBigQueryOperator(
        task_id='load_invoices_to_bq',
        bucket=GCS_BUCKET_NAME,
        source_objects=["{{ ti.xcom_pull(task_ids='extract_invoices_to_gcs', key='gcs_object_path') }}"],
        destination_project_dataset_table=f"{BIGQUERY_PROJECT_ID}.{BIGQUERY_DATASET_NAME}.invoices",
        source_format='CSV',
        skip_leading_rows=1,
        write_disposition='WRITE_APPEND',
        autodetect=True,
        gcp_conn_id=GCP_CONN_ID,
    )

    extract_tickets_to_gcs_task >> load_tickets_to_bq_task
    extract_invoices_to_gcs_task >> load_invoices_to_bq_task