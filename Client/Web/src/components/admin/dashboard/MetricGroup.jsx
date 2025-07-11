import React from 'react';
// Import SVG icons nếu cần
// import { CustomersIcon, OrdersIcon } from '../../../Svgs';

const MetricGroup = () => {
    return (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:gap-6">
            <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-gray-900 md:p-6">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gray-100 dark:bg-gray-800">
                    {/* <CustomersIcon /> */}
                </div>
                <div className="mt-5 flex items-end justify-between">
                    <div>
                        <span className="text-sm text-gray-500 dark:text-gray-400">Customers</span>
                        <h4 className="mt-2 text-2xl font-bold text-gray-800 dark:text-white/90">3,782</h4>
                    </div>
                    <span className="flex items-center gap-1 rounded-full bg-green-100 py-0.5 px-2.5 text-sm font-medium text-green-600 dark:bg-green-800/20 dark:text-green-400">
                        {/* Up Arrow SVG */}
                        11.01%
                    </span>
                </div>
            </div>
            <div className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-gray-900 md:p-6">
                <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gray-100 dark:bg-gray-800">
                    {/* <OrdersIcon /> */}
                </div>
                <div className="mt-5 flex items-end justify-between">
                    <div>
                        <span className="text-sm text-gray-500 dark:text-gray-400">Orders</span>
                        <h4 className="mt-2 text-2xl font-bold text-gray-800 dark:text-white/90">5,359</h4>
                    </div>
                    <span className="flex items-center gap-1 rounded-full bg-red-100 py-0.5 px-2.5 text-sm font-medium text-red-600 dark:bg-red-800/20 dark:text-red-400">
                        {/* Down Arrow SVG */}
                        9.05%
                    </span>
                </div>
            </div>
        </div>
    );
};

export default MetricGroup;