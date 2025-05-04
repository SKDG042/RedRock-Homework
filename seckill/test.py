#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
秒杀系统自动化压力测试脚本 - 交互式菜单版
提供简单的菜单界面，无需手动输入复杂参数
"""

import argparse
import json
import time
import random
import requests
import threading
import concurrent.futures
import csv
import os
import sys
import logging
import matplotlib.pyplot as plt
from datetime import datetime, timedelta
from tqdm import tqdm
import multiprocessing
import signal
import atexit
import inquirer  # 用于交互式菜单

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler("seckill_test.log"),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger("seckill_test")

# 统计数据结构
class Statistics:
    def __init__(self):
        self.total_requests = 0
        self.success_requests = 0
        self.failed_requests = 0
        self.total_response_time = 0
        self.min_response_time = float('inf')
        self.max_response_time = 0
        self.response_times = []
        self.status_code_count = {}
        self.error_messages = {}
        self.success_orders = []
        self.request_timestamps = []  # 记录每个请求的时间戳
        self.lock = threading.Lock()
    
    def add_result(self, success, response_time, status_code, error_msg='', order_sn=None):
        with self.lock:
            self.total_requests += 1
            self.total_response_time += response_time
            self.response_times.append(response_time)
            self.request_timestamps.append(time.time())
            
            if response_time < self.min_response_time:
                self.min_response_time = response_time
            if response_time > self.max_response_time:
                self.max_response_time = response_time
            
            if status_code in self.status_code_count:
                self.status_code_count[status_code] += 1
            else:
                self.status_code_count[status_code] = 1
            
            if success:
                self.success_requests += 1
                if order_sn:
                    self.success_orders.append(order_sn)
            else:
                self.failed_requests += 1
                if error_msg:
                    if error_msg in self.error_messages:
                        self.error_messages[error_msg] += 1
                    else:
                        self.error_messages[error_msg] = 1
    
    def get_success_rate(self):
        if self.total_requests > 0:
            return self.success_requests / self.total_requests * 100
        return 0
    
    def get_avg_response_time(self):
        if self.total_requests > 0:
            return self.total_response_time / self.total_requests
        return 0
    
    def get_percentile(self, percentile):
        if not self.response_times:
            return 0
        sorted_times = sorted(self.response_times)
        index = int(len(sorted_times) * percentile / 100)
        return sorted_times[index]
    
    def calculate_qps_over_time(self, window_size=1):
        """计算每秒请求数随时间的变化"""
        if not self.request_timestamps:
            return [], []
            
        timestamps = sorted(self.request_timestamps)
        start_time = timestamps[0]
        end_time = timestamps[-1]
        
        if end_time - start_time < window_size:
            window_size = (end_time - start_time) / 2 if end_time > start_time else 0.1
        
        time_windows = []
        current = start_time
        while current < end_time:
            time_windows.append((current, current + window_size))
            current += window_size
        
        qps_data = []
        time_points = []
        
        for window_start, window_end in time_windows:
            count = sum(1 for t in timestamps if window_start <= t < window_end)
            qps = count / window_size
            qps_data.append(qps)
            time_points.append(window_start - start_time)  # 相对时间
        
        return time_points, qps_data
    
    def save_to_csv(self, filename):
        """将统计结果保存到CSV文件"""
        with open(filename, 'w', newline='') as csvfile:
            writer = csv.writer(csvfile)
            writer.writerow(['指标', '值'])
            writer.writerow(['总请求数', self.total_requests])
            writer.writerow(['成功请求数', self.success_requests])
            writer.writerow(['失败请求数', self.failed_requests])
            writer.writerow(['成功率', f"{self.get_success_rate():.2f}%"])
            writer.writerow(['平均响应时间(ms)', f"{self.get_avg_response_time():.2f}"])
            writer.writerow(['最小响应时间(ms)', f"{self.min_response_time if self.min_response_time != float('inf') else 0:.2f}"])
            writer.writerow(['最大响应时间(ms)', f"{self.max_response_time:.2f}"])
            writer.writerow(['P50响应时间(ms)', f"{self.get_percentile(50):.2f}"])
            writer.writerow(['P90响应时间(ms)', f"{self.get_percentile(90):.2f}"])
            writer.writerow(['P95响应时间(ms)', f"{self.get_percentile(95):.2f}"])
            writer.writerow(['P99响应时间(ms)', f"{self.get_percentile(99):.2f}"])
            
            writer.writerow([])
            writer.writerow(['状态码', '数量', '百分比'])
            for code, count in self.status_code_count.items():
                percentage = count / self.total_requests * 100 if self.total_requests > 0 else 0
                writer.writerow([code, count, f"{percentage:.2f}%"])
            
            writer.writerow([])
            writer.writerow(['错误信息', '数量'])
            for msg, count in sorted(self.error_messages.items(), key=lambda x: x[1], reverse=True):
                writer.writerow([msg, count])
                
            if self.success_orders:
                writer.writerow([])
                writer.writerow(['成功订单数', len(self.success_orders)])
    
    def generate_charts(self, folder):
        """生成图表并保存"""
        os.makedirs(folder, exist_ok=True)
        
        # 1. 响应时间分布图
        if self.response_times:
            plt.figure(figsize=(10, 6))
            plt.hist(self.response_times, bins=50, alpha=0.75)
            plt.title('响应时间分布')
            plt.xlabel('响应时间 (ms)')
            plt.ylabel('请求数')
            plt.grid(True, alpha=0.3)
            plt.savefig(f"{folder}/response_time_distribution.png")
            plt.close()
        
        # 2. QPS随时间变化图
        time_points, qps_data = self.calculate_qps_over_time()
        if time_points and qps_data:
            plt.figure(figsize=(12, 6))
            plt.plot(time_points, qps_data)
            plt.title('QPS随时间变化')
            plt.xlabel('时间 (秒)')
            plt.ylabel('每秒请求数')
            plt.grid(True, alpha=0.3)
            plt.savefig(f"{folder}/qps_over_time.png")
            plt.close()
        
        # 3. 状态码分布饼图
        if self.status_code_count:
            labels = [f"HTTP {code}" for code in self.status_code_count.keys()]
            sizes = list(self.status_code_count.values())
            plt.figure(figsize=(8, 8))
            plt.pie(sizes, labels=labels, autopct='%1.1f%%', startangle=140)
            plt.axis('equal')
            plt.title('HTTP状态码分布')
            plt.savefig(f"{folder}/status_code_distribution.png")
            plt.close()
        
        # 4. 成功率与失败率饼图
        if self.total_requests > 0:
            labels = ['成功', '失败']
            sizes = [self.success_requests, self.failed_requests]
            plt.figure(figsize=(8, 8))
            plt.pie(sizes, labels=labels, autopct='%1.1f%%', startangle=140, colors=['#4CAF50', '#F44336'])
            plt.axis('equal')
            plt.title('请求成功率与失败率')
            plt.savefig(f"{folder}/success_failure_rate.png")
            plt.close()
    
    def print_results(self):
        print("\n========== 测试结果 ==========")
        print(f"总请求数: {self.total_requests}")
        
        if self.total_requests > 0:  # 避免除零错误
            success_rate = self.success_requests / self.total_requests * 100
            failed_rate = self.failed_requests / self.total_requests * 100
            print(f"成功请求数: {self.success_requests} ({success_rate:.2f}%)")
            print(f"失败请求数: {self.failed_requests} ({failed_rate:.2f}%)")
            
            avg_resp_time = self.total_response_time / self.total_requests
            print(f"平均响应时间: {avg_resp_time:.2f}ms")
            
            if self.min_response_time != float('inf'):
                print(f"最小响应时间: {self.min_response_time:.2f}ms")
                print(f"最大响应时间: {self.max_response_time:.2f}ms")
            
            # 计算响应时间分布
            if len(self.response_times) > 0:
                sorted_times = sorted(self.response_times)
                p50 = sorted_times[int(len(sorted_times) * 0.5)]
                p90 = sorted_times[int(len(sorted_times) * 0.9)]
                p95 = sorted_times[int(len(sorted_times) * 0.95)]
                p99 = sorted_times[int(len(sorted_times) * 0.99)]
                print(f"响应时间分布: P50={p50:.2f}ms, P90={p90:.2f}ms, P95={p95:.2f}ms, P99={p99:.2f}ms")
        else:
            print("没有收集到响应数据")
        
        print("\n状态码统计:")
        for code, count in self.status_code_count.items():
            percentage = count / self.total_requests * 100 if self.total_requests > 0 else 0
            print(f"  HTTP {code}: {count} ({percentage:.2f}%)")
        
        if self.error_messages:
            print("\n错误信息统计:")
            for msg, count in sorted(self.error_messages.items(), key=lambda x: x[1], reverse=True)[:5]:
                print(f"  {msg}: {count}")
            if len(self.error_messages) > 5:
                print(f"  ... 等 {len(self.error_messages)-5} 种错误")
                
        if self.success_orders:
            print(f"\n成功订单数: {len(self.success_orders)}")
            if len(self.success_orders) <= 5:
                print(f"订单号: {', '.join(self.success_orders)}")
            else:
                print(f"前5个订单号: {', '.join(self.success_orders[:5])}...")
                
        print("==============================")


class UserPool:
    """用户池，管理大量测试用户"""
    def __init__(self, base_url):
        self.base_url = base_url
        self.users = {}  # user_id -> {username, password}
        self.user_locks = {}  # user_id -> lock (防止并发使用同一用户)
        self.lock = threading.Lock()
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json'
        })
    
    def register_batch(self, count, prefix="test_user"):
        """批量注册用户"""
        logger.info(f"开始批量注册 {count} 个用户...")
        registered = 0
        batch_size = min(100, count)  # 每批次最多注册100个用户
        
        # 使用进度条
        with tqdm(total=count, desc="注册用户") as pbar:
            for i in range(0, count, batch_size):
                batch_count = min(batch_size, count - i)
                threads = []
                results = [None] * batch_count
                
                # 并发注册用户
                with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
                    for j in range(batch_count):
                        idx = i + j
                        username = f"{prefix}_{int(time.time())}_{idx}"
                        password = f"test{random.randint(100000, 999999)}"
                        
                        # 将注册任务提交到线程池
                        future = executor.submit(self._register_one, username, password, j)
                        threads.append(future)
                    
                    # 收集结果
                    for j, future in enumerate(concurrent.futures.as_completed(threads)):
                        result, idx = future.result()
                        if result:
                            results[idx] = result
                            registered += 1
                        pbar.update(1)
                
        logger.info(f"批量注册完成，成功注册 {registered}/{count} 个用户")
        return registered
    
    def _register_one(self, username, password, idx):
        """注册单个用户"""
        data = {
            "username": username,
            "password": password
        }
        
        url = f"{self.base_url}/api/user/register"
        try:
            response = self.session.post(url, json=data, timeout=10)
            
            if response.status_code == 200:
                result = response.json()
                base_resp = result.get('baseResp', {})
                if base_resp.get('code') == 0:
                    user_id = result.get('userId', 0)
                    
                    # 添加到用户池
                    with self.lock:
                        self.users[user_id] = {"username": username, "password": password}
                        self.user_locks[user_id] = threading.Lock()
                    
                    return user_id, idx
                else:
                    return None, idx
            else:
                return None, idx
        except Exception as e:
            logger.error(f"注册用户异常: {str(e)}")
            return None, idx
    
    def get_random_user(self):
        """获取随机用户ID"""
        with self.lock:
            if not self.users:
                return None
            
            # 随机选择并尝试获取锁
            user_ids = list(self.users.keys())
            random.shuffle(user_ids)
            
            for user_id in user_ids:
                lock = self.user_locks[user_id]
                if lock.acquire(blocking=False):
                    # 成功获取锁，返回用户
                    return user_id
            
            # 所有用户都在使用中
            return None
    
    def release_user(self, user_id):
        """释放用户锁"""
        if user_id in self.user_locks:
            self.user_locks[user_id].release()
    
    def get_user_count(self):
        """获取用户池中的用户数量"""
        return len(self.users)


class SeckillTest:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json'
        })
        self.stats = Statistics()
        self.user_pool = UserPool(base_url)
        self.stop_event = threading.Event()
        self.running_threads = []
    
    def create_activity(self, name, product_id, seckill_price, total_stock):
        """创建秒杀活动"""
        # 使用固定的时间戳而不是动态计算
        start_time = 1725771925  # 2024年9月8日左右
        end_time = 1825771925    # 2027年10月左右
        
        data = {
            "name": name,
            "productID": product_id,
            "seckillPrice": seckill_price,
            "startTime": start_time,
            "endTime": end_time,
            "totalStock": total_stock
        }
        
        start_time_str = datetime.fromtimestamp(start_time).strftime('%Y-%m-%d %H:%M:%S')
        end_time_str = datetime.fromtimestamp(end_time).strftime('%Y-%m-%d %H:%M:%S')
        logger.info(f"创建活动: {name}, 开始时间: {start_time_str}, 结束时间: {end_time_str}, 库存: {total_stock}")
        
        url = f"{self.base_url}/api/activity/create"
        logger.info(f"请求URL: {url}")
        start = time.time()
        try:
            response = self.session.post(url, json=data, timeout=10)
            resp_time = (time.time() - start) * 1000  # 毫秒
            
            if response.status_code == 200:
                result = response.json()
                if result.get('baseResponse', {}).get('code') == 0:
                    activity_id = result.get('activityID', 0)
                    logger.info(f"活动创建成功，ID: {activity_id}")
                    return activity_id, resp_time
                else:
                    error_msg = result.get('baseResponse', {}).get('msg', '创建活动失败')
                    logger.error(f"创建活动失败: {error_msg}")
                    return None, resp_time
            else:
                logger.error(f"创建活动请求失败: HTTP {response.status_code}, 响应: {response.text[:100]}")
                return None, resp_time
        except Exception as e:
            logger.error(f"创建活动异常: {str(e)}")
            return None, 0
    
    def get_activity(self, activity_id):
        """获取活动详情"""
        url = f"{self.base_url}/api/activity/detail/{activity_id}"
        logger.debug(f"请求URL: {url}")
        start = time.time()
        try:
            response = self.session.get(url, timeout=10)
            resp_time = (time.time() - start) * 1000  # 毫秒
            
            if response.status_code == 200:
                result = response.json()
                if result.get('baseResponse', {}).get('code') == 0:
                    activity = result.get('activity', {})
                    return activity, resp_time
                else:
                    error_msg = result.get('baseResponse', {}).get('msg', '获取活动失败')
                    logger.warning(f"获取活动失败: {error_msg}")
                    return None, resp_time
            else:
                logger.warning(f"获取活动请求失败: HTTP {response.status_code}, 响应: {response.text[:100]}")
                return None, resp_time
        except Exception as e:
            logger.error(f"获取活动异常: {str(e)}")
            return None, 0
    
    def list_activities(self):
        """获取活动列表（示例函数，根据实际API调整）"""
        url = f"{self.base_url}/api/activity/list"
        try:
            response = self.session.get(url, timeout=10)
            
            if response.status_code == 200:
                result = response.json()
                if result.get('baseResponse', {}).get('code') == 0:
                    activities = result.get('activities', [])
                    return activities
                else:
                    logger.warning("获取活动列表失败")
                    return []
            else:
                logger.warning(f"获取活动列表请求失败: HTTP {response.status_code}")
                return []
        except Exception as e:
            logger.error(f"获取活动列表异常: {str(e)}")
            return []
    
    def seckill(self, user_id, activity_id):
        """秒杀下单"""
        if not user_id:
            return False, "无有效用户ID", 0
            
        data = {
            "userID": user_id,
            "activityID": activity_id
        }
        
        url = f"{self.base_url}/api/order/seckill"
        start = time.time()
        try:
            # 添加用户ID到请求头，用于限流器识别
            headers = {'X-User-ID': str(user_id)}
            response = self.session.post(url, json=data, headers=headers, timeout=10)
            resp_time = (time.time() - start) * 1000  # 毫秒
            
            status_code = response.status_code
            try:
                result = response.json()
                base_response = result.get('baseResponse', {})
                code = base_response.get('code', -1)
                msg = base_response.get('msg', '')
                
                success = code == 0
                order_sn = None
                if success:
                    order_info = result.get('orderInfo', {})
                    order_sn = order_info.get('orderSn', '')
                
                self.stats.add_result(success, resp_time, status_code, msg if not success else '', order_sn)
                
                if success:
                    return True, order_sn, resp_time
                else:
                    return False, msg, resp_time
            except ValueError:
                # JSON解析失败，可能是非JSON响应
                error_msg = f"无法解析响应: {response.text[:100]}"
                self.stats.add_result(False, resp_time, status_code, error_msg)
                return False, error_msg, resp_time
        except Exception as e:
            error_msg = f"请求异常: {str(e)}"
            self.stats.add_result(False, 0, 0, error_msg)
            return False, error_msg, 0
    
    def continuous_seckill_worker(self, activity_id, delay_ms=0, jitter_ms=0):
        """持续进行秒杀的工作线程"""
        logger.debug(f"工作线程启动: delay={delay_ms}ms, jitter={jitter_ms}ms")
        
        while not self.stop_event.is_set():
            # 获取用户
            user_id = self.user_pool.get_random_user()
            if not user_id:
                # 用户池中没有可用用户，短暂等待后重试
                time.sleep(0.1)
                continue
            
            try:
                # 执行秒杀
                success, msg, resp_time = self.seckill(user_id, activity_id)
                
                # 计算随机延迟
                actual_delay = delay_ms
                if jitter_ms > 0:
                    actual_delay += random.randint(0, jitter_ms)
                
                # 释放用户
                self.user_pool.release_user(user_id)
                
                # 延迟一段时间
                if actual_delay > 0:
                    time.sleep(actual_delay / 1000)
            except Exception as e:
                logger.error(f"秒杀线程异常: {str(e)}")
                # 确保释放用户资源
                self.user_pool.release_user(user_id)
                time.sleep(1)  # 错误后等待一段时间再继续
    
    def run_auto_test(self, activity_id=None, product_id=1, concurrent_users=100, 
                     total_users=500, stock=1000, duration=300, 
                     delay_ms=100, jitter_ms=200, report_folder="reports"):
        """
        运行全自动化测试
        
        参数:
            activity_id: 指定活动ID，如不指定则自动创建
            product_id: 商品ID
            concurrent_users: 并发用户数
            total_users: 总用户数
            stock: 活动总库存
            duration: 测试持续时间(秒)
            delay_ms: 请求间隔延迟(毫秒)
            jitter_ms: 随机延迟抖动(毫秒)
            report_folder: 报告文件夹
        """
        start_time = datetime.now()
        test_id = f"{start_time.strftime('%Y%m%d_%H%M%S')}"
        test_report_folder = f"{report_folder}/{test_id}"
        os.makedirs(test_report_folder, exist_ok=True)
        
        # 设置信号处理，优雅停止测试
        def signal_handler(sig, frame):
            logger.info("接收到停止信号，正在优雅退出...")
            self.stop_event.set()
            sys.exit(0)
            
        signal.signal(signal.SIGINT, signal_handler)
        
        # 1. 注册用户池
        logger.info(f"开始注册用户池，目标数量: {total_users}")
        registered = self.user_pool.register_batch(total_users)
        if registered < concurrent_users:
            logger.warning(f"注册用户数({registered})小于并发用户数({concurrent_users})，将使用可用的用户数继续测试")
            concurrent_users = max(registered, 1)
        
        # 2. 创建活动或使用指定活动
        if not activity_id:
            activity_name = f"自动测试活动-{int(time.time())}"
            activity_id, _ = self.create_activity(
                name=activity_name,
                product_id=product_id,
                seckill_price=9.9,
                total_stock=stock
            )
            
            if not activity_id:
                logger.error("创建活动失败，终止测试")
                return False
            
            logger.info(f"已创建测试活动，ID: {activity_id}")
        else:
            logger.info(f"使用指定活动，ID: {activity_id}")
        
        # 3. 查询活动初始信息
        activity_before, _ = self.get_activity(activity_id)
        if activity_before:
            stock_before = activity_before.get('availableStock', 0)
            logger.info(f"活动初始库存: {stock_before}/{activity_before.get('totalStock')}")
        else:
            logger.warning("无法获取活动初始信息")
        
        # 4. 启动并发测试线程
        logger.info(f"开始启动{concurrent_users}个并发测试线程...")
        self.stop_event.clear()
        self.running_threads = []
        
        for i in range(concurrent_users):
            thread_delay = delay_ms + random.randint(0, jitter_ms // 2)  # 初始随机化避免雷同
            thread = threading.Thread(
                target=self.continuous_seckill_worker,
                args=(activity_id, thread_delay, jitter_ms),
                daemon=True
            )
            thread.start()
            self.running_threads.append(thread)
        
        # 5. 实时监控测试进度
        logger.info(f"测试开始，持续时间: {duration}秒")
        monitor_interval = 5  # 每5秒输出一次状态
        progress_bar = tqdm(total=duration, desc="测试进度")
        
        # 主循环，定期监控和更新
        elapsed = 0
        while elapsed < duration and not self.stop_event.is_set():
            time.sleep(monitor_interval)
            elapsed += monitor_interval
            progress_bar.update(monitor_interval)
            
            # 输出实时状态
            current_success = self.stats.success_requests
            current_total = self.stats.total_requests
            success_rate = (current_success / current_total * 100) if current_total > 0 else 0
            qps = current_total / elapsed if elapsed > 0 else 0
            
            logger.info(f"已处理: {current_total}请求, 成功: {current_success}({success_rate:.2f}%), " +
                        f"QPS: {qps:.2f}, 已用时间: {elapsed}/{duration}秒")
        
        progress_bar.close()
        logger.info("测试时间到，正在停止...")
        
        # 6. 停止所有测试线程
        self.stop_event.set()
        logger.info("等待所有测试线程结束...")
        
        for thread in self.running_threads:
            thread.join(timeout=2)
        
        # 7. 查询活动最终状态
        activity_after, _ = self.get_activity(activity_id)
        if activity_after and activity_before:
            final_stock = activity_after.get('availableStock', 0)
            initial_stock = activity_before.get('availableStock', 0)
            stock_reduction = initial_stock - final_stock
            logger.info(f"最终库存: {final_stock}/{activity_after.get('totalStock')}")
            logger.info(f"库存减少: {stock_reduction}")
            logger.info(f"成功订单数: {len(self.stats.success_orders)}")
            
            if stock_reduction != len(self.stats.success_orders):
                logger.warning(f"警告: 库存减少({stock_reduction})与成功订单数({len(self.stats.success_orders)})不一致!")
        
        # 8. 生成测试报告
        end_time = datetime.now()
        test_duration = (end_time - start_time).total_seconds()
        
        # 保存测试配置
        with open(f"{test_report_folder}/config.json", 'w') as f:
            json.dump({
                "test_id": test_id,
                "start_time": start_time.strftime("%Y-%m-%d %H:%M:%S"),
                "end_time": end_time.strftime("%Y-%m-%d %H:%M:%S"),
                "duration": test_duration,
                "concurrent_users": concurrent_users,
                "total_users": total_users,
                "activity_id": activity_id,
                "stock": stock,
                "delay_ms": delay_ms,
                "jitter_ms": jitter_ms,
                "base_url": self.base_url
            }, f, indent=2)
        
        # 保存统计数据
        self.stats.save_to_csv(f"{test_report_folder}/statistics.csv")
        
        # 生成图表
        self.stats.generate_charts(test_report_folder)
        
        # 生成HTML报告
        self._generate_html_report(test_report_folder, test_id, start_time, end_time)
        
        logger.info(f"测试完成! 报告已保存到: {test_report_folder}")
        self.stats.print_results()
        
        return True

    def _generate_html_report(self, folder, test_id, start_time, end_time):
        """生成HTML测试报告"""
        duration = (end_time - start_time).total_seconds()
        
        html_content = f"""
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>秒杀系统压力测试报告 - {test_id}</title>
            <style>
                body {{ font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 1200px; margin: 0 auto; padding: 20px; }}
                h1, h2, h3 {{ color: #0066cc; }}
                .header {{ background-color: #f5f5f5; padding: 20px; border-radius: 5px; margin-bottom: 20px; }}
                .summary {{ display: flex; flex-wrap: wrap; gap: 20px; margin-bottom: 30px; }}
                .summary-card {{ background-color: #f9f9f9; border-radius: 5px; padding: 15px; flex: 1; min-width: 200px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }}
                .chart-container {{ display: flex; flex-wrap: wrap; gap: 20px; margin-bottom: 30px; }}
                .chart {{ flex: 1; min-width: 45%; background-color: #fff; border-radius: 5px; padding: 15px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }}
                table {{ width: 100%; border-collapse: collapse; margin-bottom: 20px; }}
                th, td {{ padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }}
                th {{ background-color: #f2f2f2; }}
                tr:hover {{ background-color: #f5f5f5; }}
                .success {{ color: #4CAF50; }}
                .failure {{ color: #F44336; }}
            </style>
        </head>
        <body>
            <div class="header">
                <h1>秒杀系统压力测试报告</h1>
                <p>测试ID: {test_id}</p>
                <p>开始时间: {start_time.strftime("%Y-%m-%d %H:%M:%S")}</p>
                <p>结束时间: {end_time.strftime("%Y-%m-%d %H:%M:%S")}</p>
                <p>测试持续时间: {duration:.2f}秒</p>
            </div>
            
            <h2>测试结果摘要</h2>
            <div class="summary">
                <div class="summary-card">
                    <h3>请求统计</h3>
                    <p>总请求数: {self.stats.total_requests}</p>
                    <p>成功请求数: <span class="success">{self.stats.success_requests}</span></p>
                    <p>失败请求数: <span class="failure">{self.stats.failed_requests}</span></p>
                    <p>成功率: {self.stats.get_success_rate():.2f}%</p>
                </div>
                
                <div class="summary-card">
                    <h3>性能指标</h3>
                    <p>平均QPS: {self.stats.total_requests / duration:.2f}</p>
                    <p>平均响应时间: {self.stats.get_avg_response_time():.2f}ms</p>
                    <p>最小响应时间: {self.stats.min_response_time if self.stats.min_response_time != float('inf') else 0:.2f}ms</p>
                    <p>最大响应时间: {self.stats.max_response_time:.2f}ms</p>
                </div>
                
                <div class="summary-card">
                    <h3>响应时间分布</h3>
                    <p>P50: {self.stats.get_percentile(50):.2f}ms</p>
                    <p>P90: {self.stats.get_percentile(90):.2f}ms</p>
                    <p>P95: {self.stats.get_percentile(95):.2f}ms</p>
                    <p>P99: {self.stats.get_percentile(99):.2f}ms</p>
                </div>
            </div>
            
            <h2>图表分析</h2>
            <div class="chart-container">
                <div class="chart">
                    <h3>响应时间分布</h3>
                    <img src="response_time_distribution.png" alt="响应时间分布" style="width: 100%;">
                </div>
                
                <div class="chart">
                    <h3>QPS随时间变化</h3>
                    <img src="qps_over_time.png" alt="QPS随时间变化" style="width: 100%;">
                </div>
                
                <div class="chart">
                    <h3>状态码分布</h3>
                    <img src="status_code_distribution.png" alt="状态码分布" style="width: 100%;">
                </div>
                
                <div class="chart">
                    <h3>成功率与失败率</h3>
                    <img src="success_failure_rate.png" alt="成功率与失败率" style="width: 100%;">
                </div>
            </div>
            
            <h2>状态码分布</h2>
            <table>
                <tr>
                    <th>状态码</th>
                    <th>数量</th>
                    <th>百分比</th>
                </tr>
        """
        
        # 添加状态码表格数据
        for code, count in self.stats.status_code_count.items():
            percentage = count / self.stats.total_requests * 100 if self.stats.total_requests > 0 else 0
            html_content += f"""
                <tr>
                    <td>HTTP {code}</td>
                    <td>{count}</td>
                    <td>{percentage:.2f}%</td>
                </tr>
            """
        
        html_content += """
            </table>
            
            <h2>错误信息分析</h2>
            <table>
                <tr>
                    <th>错误信息</th>
                    <th>数量</th>
                    <th>百分比</th>
                </tr>
        """
        
        # 添加错误信息表格数据
        for msg, count in sorted(self.stats.error_messages.items(), key=lambda x: x[1], reverse=True):
            percentage = count / self.stats.failed_requests * 100 if self.stats.failed_requests > 0 else 0
            html_content += f"""
                <tr>
                    <td>{msg}</td>
                    <td>{count}</td>
                    <td>{percentage:.2f}%</td>
                </tr>
            """
        
        html_content += """
            </table>
            
            <h2>订单信息</h2>
            <p>成功创建订单数: """ + str(len(self.stats.success_orders)) + """</p>
            
            <script>
                // 可以在这里添加一些交互式JavaScript功能
                document.addEventListener('DOMContentLoaded', function() {
                    console.log('测试报告加载完成');
                });
            </script>
        </body>
        </html>
        """
        
        # 写入HTML文件
        with open(f"{folder}/report.html", 'w', encoding='utf-8') as f:
            f.write(html_content)


class LoadTestProfile:
    """负载测试配置文件"""
    def __init__(self, name, concurrent_users, total_users, duration, delay_ms, jitter_ms):
        self.name = name
        self.concurrent_users = concurrent_users
        self.total_users = total_users
        self.duration = duration
        self.delay_ms = delay_ms
        self.jitter_ms = jitter_ms
    
    @classmethod
    def get_predefined_profiles(cls):
        """获取预定义的测试配置"""
        return {
            "light": cls("轻量测试", 10, 50, 60, 500, 200),
            "medium": cls("中等负载", 50, 200, 120, 200, 300),
            "heavy": cls("高负载", 100, 500, 300, 100, 200),
            "extreme": cls("极限测试", 200, 1000, 600, 50, 100),
            "endurance": cls("持久测试", 50, 500, 1800, 200, 300),
        }
    
    @classmethod
    def get_profile_choices(cls):
        """获取配置选项列表，用于菜单显示"""
        profiles = cls.get_predefined_profiles()
        choices = []
        for key, profile in profiles.items():
            choices.append({
                'name': f"{profile.name} - {profile.concurrent_users}并发/{profile.duration}秒/{profile.total_users}用户",
                'value': key
            })
        return choices


# 交互式菜单处理类
class SeckillMenu:
    def __init__(self):
        self.base_url = "http://localhost:8080"
        self.test = SeckillTest(self.base_url)
        self.report_folder = "reports"
        os.makedirs(self.report_folder, exist_ok=True)
    
    def show_main_menu(self):
        """显示主菜单"""
        choices = [
            {
                'type': 'list',
                'name': 'action',
                'message': '请选择要执行的操作:',
                'choices': [
                    {'name': '1. 运行自动化测试', 'value': 'auto_test'},
                    {'name': '2. 创建秒杀活动', 'value': 'create_activity'},
                    {'name': '3. 批量注册用户', 'value': 'register_users'},
                    {'name': '4. 查询活动信息', 'value': 'query_activity'},
                    {'name': '5. 修改服务器地址', 'value': 'change_server'},
                    {'name': '6. 退出程序', 'value': 'exit'}
                ]
            }
        ]
        
        answers = inquirer.prompt(choices)
        action = answers['action']
        
        if action == 'auto_test':
            self.run_auto_test_menu()
        elif action == 'create_activity':
            self.create_activity_menu()
        elif action == 'register_users':
            self.register_users_menu()
        elif action == 'query_activity':
            self.query_activity_menu()
        elif action == 'change_server':
            self.change_server_menu()
        elif action == 'exit':
            print("程序已退出")
            sys.exit(0)
    
    def run_auto_test_menu(self):
        """自动化测试菜单"""
        # 1. 选择测试配置
        profile_choices = LoadTestProfile.get_profile_choices()
        profile_choices.append({'name': '自定义配置', 'value': 'custom'})
        
        questions = [
            {
                'type': 'list',
                'name': 'profile',
                'message': '请选择测试配置:',
                'choices': profile_choices
            }
        ]
        
        answers = inquirer.prompt(questions)
        profile_key = answers['profile']
        
        # 处理自定义配置
        if profile_key == 'custom':
            custom_questions = [
                {
                    'type': 'input',
                    'name': 'concurrent_users',
                    'message': '并发用户数:',
                    'default': '50',
                    'validate': lambda x: x.isdigit() and int(x) > 0
                },
                {
                    'type': 'input',
                    'name': 'total_users',
                    'message': '总用户数(建议比并发用户数大):',
                    'default': '200',
                    'validate': lambda x: x.isdigit() and int(x) > 0
                },
                {
                    'type': 'input',
                    'name': 'duration',
                    'message': '测试持续时间(秒):',
                    'default': '300',
                    'validate': lambda x: x.isdigit() and int(x) > 0
                },
                {
                    'type': 'input',
                    'name': 'delay_ms',
                    'message': '请求间隔延迟(毫秒):',
                    'default': '100',
                    'validate': lambda x: x.isdigit() and int(x) >= 0
                },
                {
                    'type': 'input',
                    'name': 'jitter_ms',
                    'message': '随机延迟抖动(毫秒):',
                    'default': '200',
                    'validate': lambda x: x.isdigit() and int(x) >= 0
                }
            ]
            
            custom_answers = inquirer.prompt(custom_questions)
            concurrent_users = int(custom_answers['concurrent_users'])
            total_users = int(custom_answers['total_users'])
            duration = int(custom_answers['duration'])
            delay_ms = int(custom_answers['delay_ms'])
            jitter_ms = int(custom_answers['jitter_ms'])
        else:
            # 使用预定义配置
            profiles = LoadTestProfile.get_predefined_profiles()
            profile = profiles[profile_key]
            concurrent_users = profile.concurrent_users
            total_users = profile.total_users
            duration = profile.duration
            delay_ms = profile.delay_ms
            jitter_ms = profile.jitter_ms
            
            print(f"使用预设配置: {profile.name}")
            print(f"- 并发用户: {concurrent_users}")
            print(f"- 总用户数: {total_users}")
            print(f"- 持续时间: {duration}秒")
            print(f"- 请求延迟: {delay_ms}ms (抖动: {jitter_ms}ms)")
        
        # 2. 选择活动创建方式
        activity_questions = [
            {
                'type': 'list',
                'name': 'activity_option',
                'message': '请选择秒杀活动:',
                'choices': [
                    {'name': '创建新活动', 'value': 'new'},
                    {'name': '使用现有活动', 'value': 'existing'}
                ]
            }
        ]
        
        activity_answers = inquirer.prompt(activity_questions)
        
        if activity_answers['activity_option'] == 'new':
            # 创建新活动
            activity_params = [
                {
                    'type': 'input',
                    'name': 'activity_name',
                    'message': '活动名称:',
                    'default': f'自动测试活动-{int(time.time())}',
                },
                {
                    'type': 'input',
                    'name': 'product_id',
                    'message': '商品ID:',
                    'default': '1',
                    'validate': lambda x: x.isdigit() and int(x) > 0
                },
                {
                    'type': 'input',
                    'name': 'seckill_price',
                    'message': '秒杀价格:',
                    'default': '9.9',
                    'validate': lambda x: float(x) > 0
                },
                {
                    'type': 'input',
                    'name': 'stock',
                    'message': '库存数量:',
                    'default': '1000',
                    'validate': lambda x: x.isdigit() and int(x) > 0
                }
            ]
            
            activity_param_answers = inquirer.prompt(activity_params)
            
            print("正在创建活动...")
            activity_id, _ = self.test.create_activity(
                name=activity_param_answers['activity_name'],
                product_id=int(activity_param_answers['product_id']),
                seckill_price=float(activity_param_answers['seckill_price']),
                total_stock=int(activity_param_answers['stock'])
            )
            
            if not activity_id:
                print("创建活动失败，返回主菜单")
                return
                
            print(f"活动创建成功，ID: {activity_id}")
            stock = int(activity_param_answers['stock'])
                
        else:
            # 使用现有活动
            activity_id_question = [
                {
                    'type': 'input',
                    'name': 'activity_id',
                    'message': '请输入现有活动ID:',
                    'validate': lambda x: x.isdigit() and int(x) > 0
                }
            ]
            
            activity_id_answer = inquirer.prompt(activity_id_question)
            activity_id = int(activity_id_answer['activity_id'])
            
            # 查询活动库存
            activity, _ = self.test.get_activity(activity_id)
            if not activity:
                print(f"无法获取活动ID {activity_id} 的信息，返回主菜单")
                return
                
            stock = activity.get('totalStock', 1000)
            print(f"活动 '{activity.get('name')}' 当前库存: {activity.get('availableStock')}/{stock}")
        
        # 3. 确认启动测试
        confirm_questions = [
            {
                'type': 'confirm',
                'name': 'confirm',
                'message': '确认开始测试?',
                'default': True
            }
        ]
        
        confirm_answers = inquirer.prompt(confirm_questions)
        
        if not confirm_answers['confirm']:
            print("已取消测试，返回主菜单")
            return
        
        # 4. 启动测试
        print("\n========== 开始测试 ==========")
        print(f"活动ID: {activity_id}")
        print(f"并发用户: {concurrent_users}, 总用户: {total_users}")
        print(f"持续时间: {duration}秒")
        print(f"请求间隔: {delay_ms}ms, 抖动: {jitter_ms}ms")
        print("==============================\n")
        
        self.test.run_auto_test(
            activity_id=activity_id,
            concurrent_users=concurrent_users,
            total_users=total_users,
            stock=stock,
            duration=duration,
            delay_ms=delay_ms,
            jitter_ms=jitter_ms,
            report_folder=self.report_folder
        )
        
        # 测试完成，等待用户返回主菜单
        input("\n按Enter键返回主菜单...")
    
    def create_activity_menu(self):
        """创建活动菜单"""
        questions = [
            {
                'type': 'input',
                'name': 'activity_name',
                'message': '活动名称:',
                'default': f'测试活动-{int(time.time())}',
            },
            {
                'type': 'input',
                'name': 'product_id',
                'message': '商品ID:',
                'default': '1',
                'validate': lambda x: x.isdigit() and int(x) > 0
            },
            {
                'type': 'input',
                'name': 'seckill_price',
                'message': '秒杀价格:',
                'default': '9.9',
                'validate': lambda x: float(x) > 0
            },
            {
                'type': 'input',
                'name': 'stock',
                'message': '库存数量:',
                'default': '1000',
                'validate': lambda x: x.isdigit() and int(x) > 0
            }
        ]
        
        answers = inquirer.prompt(questions)
        
        print("正在创建活动...")
        activity_id, _ = self.test.create_activity(
            name=answers['activity_name'],
            product_id=int(answers['product_id']),
            seckill_price=float(answers['seckill_price']),
            total_stock=int(answers['stock'])
        )
        
        if activity_id:
            print(f"活动创建成功，ID: {activity_id}")
            
            # 查询活动详情
            activity, _ = self.test.get_activity(activity_id)
            if activity:
                print(f"活动状态: {activity.get('status')}")
                print(f"开始时间: {datetime.fromtimestamp(activity.get('startTime', 0)).strftime('%Y-%m-%d %H:%M:%S')}")
                print(f"结束时间: {datetime.fromtimestamp(activity.get('endTime', 0)).strftime('%Y-%m-%d %H:%M:%S')}")
        else:
            print("活动创建失败")
        
        input("\n按Enter键返回主菜单...")
    
    def register_users_menu(self):
        """批量注册用户菜单"""
        questions = [
            {
                'type': 'input',
                'name': 'user_count',
                'message': '要注册的用户数量:',
                'default': '100',
                'validate': lambda x: x.isdigit() and int(x) > 0
            },
            {
                'type': 'input',
                'name': 'prefix',
                'message': '用户名前缀:',
                'default': 'test_user',
            }
        ]
        
        answers = inquirer.prompt(questions)
        user_count = int(answers['user_count'])
        prefix = answers['prefix']
        
        print(f"开始注册 {user_count} 个用户...")
        registered = self.test.user_pool.register_batch(user_count, prefix)
        
        print(f"注册完成! 成功注册 {registered}/{user_count} 个用户")
        input("\n按Enter键返回主菜单...")
    
    def query_activity_menu(self):
        """查询活动菜单"""
        questions = [
            {
                'type': 'input',
                'name': 'activity_id',
                'message': '活动ID:',
                'validate': lambda x: x.isdigit() and int(x) > 0
            }
        ]
        
        answers = inquirer.prompt(questions)
        activity_id = int(answers['activity_id'])
        
        print(f"正在查询活动ID {activity_id} 的信息...")
        activity, _ = self.test.get_activity(activity_id)
        
        if activity:
            print("\n活动详情:")
            print(f"ID: {activity.get('id')}")
            print(f"名称: {activity.get('name')}")
            print(f"商品ID: {activity.get('productId')}")
            print(f"秒杀价格: {activity.get('seckillPrice')}")
            print(f"原价: {activity.get('originalPrice')}")
            print(f"当前库存: {activity.get('availableStock')}/{activity.get('totalStock')}")
            print(f"活动状态: {activity.get('status')}")
            print(f"开始时间: {datetime.fromtimestamp(activity.get('startTime', 0)).strftime('%Y-%m-%d %H:%M:%S')}")
            print(f"结束时间: {datetime.fromtimestamp(activity.get('endTime', 0)).strftime('%Y-%m-%d %H:%M:%S')}")
        else:
            print(f"未找到ID为 {activity_id} 的活动")
        
        input("\n按Enter键返回主菜单...")
    
    def change_server_menu(self):
        """修改服务器地址菜单"""
        questions = [
            {
                'type': 'input',
                'name': 'server_url',
                'message': '服务器地址:',
                'default': self.base_url,
            }
        ]
        
        answers = inquirer.prompt(questions)
        new_url = answers['server_url']
        
        if new_url != self.base_url:
            self.base_url = new_url
            self.test = SeckillTest(self.base_url)
            print(f"服务器地址已更新为: {self.base_url}")
        else:
            print("服务器地址未变更")
        
        input("\n按Enter键返回主菜单...")
    
    def run(self):
        """运行菜单主循环"""
        print("\n欢迎使用秒杀系统测试工具")
        print(f"当前服务器: {self.base_url}")
        
        try:
            while True:
                self.show_main_menu()
        except KeyboardInterrupt:
            print("\n程序已退出")
            sys.exit(0)


def main():
    # 检查依赖包
    try:
        import inquirer
    except ImportError:
        print("缺少必要的依赖包。正在安装...")
        import subprocess
        subprocess.check_call([sys.executable, "-m", "pip", "install", "inquirer", "matplotlib", "tqdm"])
        print("依赖包安装完成，请重新运行程序")
        sys.exit(0)
    
    # 启动交互式菜单
    menu = SeckillMenu()
    menu.run()


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n程序已退出")
    except Exception as e:
        logger.exception(f"程序执行异常: {str(e)}")