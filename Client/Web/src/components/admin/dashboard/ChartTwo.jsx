import React, { useEffect } from 'react';
// import ApexCharts from 'apexcharts';

const ChartTwo = () => {
    useEffect(() => {
       // Logic để render biểu đồ tròn
    }, []);

    return (
        <div className="rounded-2xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900">
            <div className="px-5 pt-5 sm:px-6 sm:pt-6">
                <div className="flex justify-between">
                    <div>
                        <h3 className="text-lg font-semibold text-gray-800 dark:text-white/90">Monthly Target</h3>
                        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">Target you’ve set for each month</p>
                    </div>
                </div>
                <div className="relative my-4 h-[195px]">
                    <div id="chartTwo" className="h-full"></div>
                    <span className="absolute left-1/2 top-[85%] -translate-x-1/2 -translate-y-[85%] rounded-full bg-green-100 px-3 py-1 text-xs font-medium text-green-600 dark:bg-green-800/20 dark:text-green-400">+10%</span>
                </div>
                 <p className="mx-auto mt-1.5 w-full max-w-[380px] text-center text-sm text-gray-500">
                    You earn $3287 today, it's higher than last month. Keep up your good work!
                </p>
            </div>
             <div className="flex items-center justify-center gap-5 border-t border-gray-200 dark:border-gray-800 px-6 py-3.5 sm:gap-8 sm:py-5 mt-4">
                {/* Target, Revenue, Today sections */}
            </div>
        </div>
    );
};
export default ChartTwo;