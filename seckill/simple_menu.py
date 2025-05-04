#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import sys
import os
import time
from datetime import datetime
import concurrent.futures
import random

# 导入原测试脚本的关键类
from test import SeckillTest, LoadTestProfile, Statistics, UserPool

def simple_menu():
    test = SeckillTest("http://localhost:8080")
    report_folder = "reports"
    os.makedirs(report_folder, exist_ok=True)
    
    while True:
        print("\n===== 秒杀系统测试工具 =====")
        print("1. 运行轻量测试 (10并发/60秒)")
        print("2. 运行中等负载 (50并发/120秒)")
        print("3. 运行高负载 (100并发/300秒)")
        print("4. 极限测试 (200并发/600秒)")
        print("5. 创建新活动")
        print("6. 批量注册用户")
        print("7. 查询活动信息")
        print("8. 验证系统质量")
        print("9. 修改服务器地址")
        print("10. 退出程序")
        
        choice = input("\n请选择: ")
        
        if choice == "1":
            run_test(test, "light", report_folder)
        elif choice == "2":
            run_test(test, "medium", report_folder)
        elif choice == "3":
            run_test(test, "heavy", report_folder)
        elif choice == "4":
            run_test(test, "extreme", report_folder)
        elif choice == "5":
            create_activity(test)
        elif choice == "6":
            register_users(test)
        elif choice == "7":
            query_activity(test)
        elif choice == "8":
            verify_menu(test)
        elif choice == "9":
            change_server(test)
        elif choice == "10":
            print("程序已退出")
            break
        else:
            print("无效选择，请重试")

def run_test(test, profile_key, report_folder):
    try:
        profiles = LoadTestProfile.get_predefined_profiles()
        profile = profiles[profile_key]
        
        print(f"\n使用配置: {profile.name}")
        print(f"- 并发用户: {profile.concurrent_users}")
        print(f"- 总用户数: {profile.total_users}")
        print(f"- 持续时间: {profile.duration}秒")
        print(f"- 请求延迟: {profile.delay_ms}ms (抖动: {profile.jitter_ms}ms)")
        
        option = input("\n选择活动: (1) 创建新活动 (2) 使用现有活动: ")
        
        if option == "1":
            name = f"自动测试活动-{int(time.time())}"
            product_id = 1
            try:
                product_id = int(input("商品ID [1]: ") or "1")
            except:
                pass
                
            price = 9.9
            try:
                price = float(input("秒杀价格 [9.9]: ") or "9.9")
            except:
                pass
                
            stock = 1000
            try:
                stock = int(input("库存数量 [1000]: ") or "1000")
            except:
                pass
                
            print("正在创建活动...")
            activity_id, _ = test.create_activity(name, product_id, price, stock)
            
            if not activity_id:
                print("创建活动失败")
                return
                
            print(f"活动创建成功，ID: {activity_id}")
        else:
            try:
                activity_id = int(input("请输入活动ID: "))
                
                # 查询活动库存
                activity, _ = test.get_activity(activity_id)
                if not activity:
                    print(f"无法获取活动ID {activity_id} 的信息")
                    return
                    
                stock = activity.get('totalStock', 1000)
                print(f"活动 '{activity.get('name')}' 当前库存: {activity.get('availableStock')}/{stock}")
            except:
                print("无效的活动ID")
                return
        
        confirm = input("\n确认开始测试? (y/n): ")
        if confirm.lower() != 'y':
            print("已取消测试")
            return
            
        test.run_auto_test(
            activity_id=activity_id,
            concurrent_users=profile.concurrent_users,
            total_users=profile.total_users,
            duration=profile.duration,
            delay_ms=profile.delay_ms,
            jitter_ms=profile.jitter_ms,
            report_folder=report_folder
        )
    except Exception as e:
        print(f"测试出错: {str(e)}")

def create_activity(test):
    try:
        name = input("活动名称 [测试活动]: ") or f"测试活动-{int(time.time())}"
        
        product_id = 1
        try:
            product_id = int(input("商品ID [1]: ") or "1")
        except:
            pass
            
        price = 9.9
        try:
            price = float(input("秒杀价格 [9.9]: ") or "9.9")
        except:
            pass
            
        stock = 1000
        try:
            stock = int(input("库存数量 [1000]: ") or "1000")
        except:
            pass
        
        print("正在创建活动...")
        activity_id, _ = test.create_activity(name, product_id, price, stock)
        
        if activity_id:
            print(f"活动创建成功，ID: {activity_id}")
            
            # 查询活动详情
            activity, _ = test.get_activity(activity_id)
            if activity:
                print(f"活动状态: {activity.get('status')}")
                print(f"开始时间: {datetime.fromtimestamp(activity.get('startTime', 0)).strftime('%Y-%m-%d %H:%M:%S')}")
                print(f"结束时间: {datetime.fromtimestamp(activity.get('endTime', 0)).strftime('%Y-%m-%d %H:%M:%S')}")
        else:
            print("活动创建失败")
    except Exception as e:
        print(f"创建活动出错: {str(e)}")

def register_users(test):
    try:
        count = 100
        try:
            count = int(input("注册用户数量 [100]: ") or "100")
        except:
            pass
            
        prefix = input("用户名前缀 [test_user]: ") or "test_user"
        
        print(f"开始注册 {count} 个用户...")
        registered = test.user_pool.register_batch(count, prefix)
        
        print(f"注册完成! 成功: {registered}/{count} 个用户")
    except Exception as e:
        print(f"注册用户出错: {str(e)}")

def query_activity(test):
    try:
        activity_id = int(input("活动ID: "))
        
        print(f"正在查询活动ID {activity_id}...")
        activity, _ = test.get_activity(activity_id)
        
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
    except Exception as e:
        print(f"查询活动出错: {str(e)}")

def verify_menu(test):
    print("\n===== 系统质量验证 =====")
    print("1. 检查超卖问题")
    print("2. 验证同用户重复购买")
    print("3. 测试接口幂等性(单用户)")
    print("4. 测试接口幂等性(并发模式)")
    print("5. 返回主菜单")
    
    choice = input("\n请选择: ")
    
    if choice == "1":
        check_oversold(test)
    elif choice == "2":
        verify_duplicate_purchase(test)
    elif choice == "3":
        verify_idempotency(test)
    elif choice == "4":
        verify_idempotency_concurrent(test)
    elif choice == "5":
        return
    else:
        print("无效选择")

def check_oversold(test):
    try:
        activity_id = int(input("请输入活动ID: "))
        
        # 查询活动
        activity, _ = test.get_activity(activity_id)
        if not activity:
            print(f"未找到ID为 {activity_id} 的活动")
            return
        
        init_stock = activity.get('totalStock', 0)
        avail_stock = activity.get('availableStock', 0)
        used_stock = init_stock - avail_stock
        
        print(f"\n活动 '{activity.get('name')}' 库存状态:")
        print(f"初始库存: {init_stock}")
        print(f"剩余库存: {avail_stock}")
        print(f"消耗库存: {used_stock}")
        
        # 高并发测试选项
        print("\n===== 极限并发测试选项 =====")
        print("1. 高并发 (500并发/1000用户)")
        print("2. 超高并发 (1000并发/2000用户)")
        print("3. 极限并发 (2000并发/5000用户)")
        print("4. 自定义并发")
        print("5. 返回")
        
        option = input("\n选择并发级别: ")
        
        if option == "5":
            return
            
        if option == "1":
            concurrent_users = 500
            total_users = 1000
            duration = 30
        elif option == "2":
            concurrent_users = 1000
            total_users = 2000
            duration = 60
        elif option == "3":
            concurrent_users = 2000
            total_users = 5000
            duration = 90
        elif option == "4":
            try:
                concurrent_users = int(input("并发用户数: "))
                total_users = int(input("总用户数: "))
                duration = int(input("测试持续时间(秒): "))
                
                if total_users < concurrent_users:
                    total_users = concurrent_users * 2
                    print(f"总用户数已调整为: {total_users}")
            except ValueError:
                print("输入无效，使用默认值")
                concurrent_users = 500
                total_users = 1000
                duration = 30
        else:
            print("选择无效，使用高并发设置")
            concurrent_users = 500
            total_users = 1000
            duration = 30
                
        if avail_stock < 10:
            print("剩余库存太少，无法进行有效测试")
            return
            
        print(f"\n将执行极限并发测试:")
        print(f"- 并发用户数: {concurrent_users}")
        print(f"- 测试用户总数: {total_users}")
        print(f"- 测试持续时间: {duration}秒")
        print(f"- 测试前剩余库存: {avail_stock}")
            
        confirm = input("\n确认执行测试? (y/n): ")
        if confirm.lower() != 'y':
            print("已取消测试")
            return
            
        # 注册测试用户
        print(f"注册{total_users}个测试用户...")
        registered = test.user_pool.register_batch(total_users, f"stress_{int(time.time())}")
        print(f"成功注册 {registered} 个用户")
        
        if registered < concurrent_users * 0.5:
            print(f"警告: 注册用户数不足，只有{registered}个，可能影响测试结果")
            confirm = input("是否继续? (y/n): ")
            if confirm.lower() != 'y':
                return
        
        # 清空统计信息
        test.stats = Statistics()
        
        # 保存原始库存
        orig_stock = avail_stock
        
        # 执行直接并发请求，不生成报告和图表
        print("\n开始执行极限并发测试，请稍候...")
        
        def simple_seckill_worker(user_id, activity_id):
            return test.seckill(user_id, activity_id)
            
        # 获取用户ID列表
        user_ids = list(test.user_pool.users.keys())[:registered]
        
        success_count = 0
        fail_count = 0
        start_time = time.time()
        
        # 直接使用多进程/线程池进行并发请求
        with concurrent.futures.ThreadPoolExecutor(max_workers=concurrent_users) as executor:
            futures = []
            for _ in range(min(len(user_ids), concurrent_users * 10)):  # 增加请求次数
                user_id = user_ids[_ % len(user_ids)]  # 循环使用用户
                futures.append(executor.submit(simple_seckill_worker, user_id, activity_id))
            
            # 等待所有任务完成或达到时间限制
            done, not_done = concurrent.futures.wait(
                futures, 
                timeout=duration,
                return_when=concurrent.futures.FIRST_EXCEPTION
            )
            
            # 取消未完成的任务
            for future in not_done:
                future.cancel()
                
            # 处理已完成的任务结果
            for future in done:
                try:
                    success, msg, _ = future.result()
                    if success:
                        success_count += 1
                    else:
                        fail_count += 1
                except Exception:
                    fail_count += 1
        
        end_time = time.time()
        actual_duration = end_time - start_time
        
        # 验证结果
        activity_after, _ = test.get_activity(activity_id)
        final_stock = activity_after.get('availableStock', 0)
        stock_reduction = orig_stock - final_stock
        
        print("\n===== 极限并发测试结果 =====")
        print(f"测试持续时间: {actual_duration:.2f}秒")
        print(f"总请求数: {len(done)}")
        print(f"成功请求数: {success_count}")
        print(f"失败请求数: {fail_count}")
        print(f"QPS: {len(done)/actual_duration:.2f}次/秒")
        print(f"\n初始库存: {orig_stock}")
        print(f"测试后库存: {final_stock}")
        print(f"库存减少量: {stock_reduction}")
        
        if stock_reduction == success_count:
            print("\n✅ 库存减少量与成功订单数一致，未发现超卖")
        else:
            print(f"\n⚠️ 警告: 库存减少量({stock_reduction})与成功订单数({success_count})不一致!")
            if stock_reduction > success_count:
                print("   可能存在未完成订单的库存锁定")
            elif stock_reduction < success_count:
                print("   检测到超卖情况! 系统存在严重安全隐患!")
        
        # 整体超卖风险评估
        if used_stock + stock_reduction > init_stock:
            print("\n⚠️ 警告: 当前活动已存在超卖情况!")
        else:
            print("\n当前活动整体库存状态正常")
            
    except Exception as e:
        print(f"检查超卖出错: {str(e)}")

def verify_duplicate_purchase(test):
    try:
        activity_id = int(input("请输入活动ID: "))
        
        # 查询活动
        activity, _ = test.get_activity(activity_id)
        if not activity:
            print(f"未找到ID为 {activity_id} 的活动")
            return
            
        if activity.get('availableStock', 0) <= 0:
            print("活动库存已售罄，无法进行测试")
            return
            
        print(f"\n验证活动 '{activity.get('name')}' 的重复购买限制")
        
        # 选择测试模式
        print("\n测试模式:")
        print("1. 单用户重复购买测试")
        print("2. 多用户并发购买测试 (100用户)")
        print("3. 返回")
        
        option = input("\n选择测试模式: ")
        
        if option == "3":
            return
            
        if option == "1":
            # 单用户测试模式
            print("\n=== 单用户重复购买测试 ===")
            # 注册测试用户
            print("注册测试用户...")
            prefix = f"dup_test_{int(time.time())}"
            registered = test.user_pool.register_batch(1, prefix)
            
            if registered <= 0:
                print("注册测试用户失败")
                return
                
            # 获取注册的用户ID - 从用户池中获取
            user_ids = list(test.user_pool.users.keys())
            if not user_ids:
                print("找不到已注册的用户")
                return
                
            # 使用最后注册的用户
            user_id = user_ids[-1]
            print(f"使用测试用户ID: {user_id}")
                
            # 第一次秒杀
            print("\n第一次秒杀请求...")
            success1, order_sn1, _ = test.seckill(user_id, activity_id)
            
            if not success1:
                print(f"首次秒杀失败，无法验证: {order_sn1}")
                return
                
            print(f"首次秒杀成功，订单号: {order_sn1}")
            
            # 第二次秒杀
            print("\n第二次秒杀请求(相同用户)...")
            success2, order_sn2, _ = test.seckill(user_id, activity_id)
            
            if success2:
                print(f"⚠️ 警告: 重复秒杀成功，订单号: {order_sn2}")
                print("系统存在重复购买风险!")
            else:
                print(f"✅ 重复秒杀被拦截: {order_sn2}")
                print("重复购买限制正常工作")
                
        elif option == "2":
            # 多用户并发测试模式
            print("\n=== 多用户并发购买测试 ===")
            
            users_count = 100
            print(f"注册{users_count}个测试用户...")
            registered = test.user_pool.register_batch(users_count, f"dup_test_multi_{int(time.time())}")
            
            if registered < 10:  # 至少需要10个用户
                print(f"注册用户不足，只有{registered}个用户")
                return
                
            print(f"成功注册{registered}个用户，开始并发测试")
            
            # 每个用户尝试购买两次
            max_attempts = 2
            print(f"每个用户将尝试购买{max_attempts}次")
            
            # 执行并发测试
            concurrent_users = min(50, registered)
            duration = 60  # 60秒
            
            print(f"启动{concurrent_users}用户并发测试，持续{duration}秒...")
            
            # 清空统计信息
            test.stats = Statistics()
            
            # 保存原始库存值
            orig_stock = activity.get('availableStock', 0)
            
            # 使用自定义测试逻辑，确保每个用户发送两次请求
            # 这里需要修改run_auto_test方法，或自定义测试逻辑
            print("开始测试，每个用户将发送多次请求...")
            
            def multi_purchase_worker(user_id, activity_id, attempts):
                results = []
                for i in range(attempts):
                    success, msg, _ = test.seckill(user_id, activity_id)
                    time.sleep(0.1)  # 短暂暂停，避免请求过于密集
                    results.append((success, msg))
                return results
            
            # 获取用户ID列表
            user_ids = list(test.user_pool.users.keys())[:registered]
            
            # 使用线程池并发执行
            with concurrent.futures.ThreadPoolExecutor(max_workers=concurrent_users) as executor:
                futures = []
                for user_id in user_ids:
                    futures.append(executor.submit(
                        multi_purchase_worker, user_id, activity_id, max_attempts
                    ))
                
                # 等待所有任务完成
                all_results = []
                for future in concurrent.futures.as_completed(futures):
                    all_results.extend(future.result())
            
            # 分析结果
            success_count = sum(1 for r in all_results if r[0])
            fail_count = len(all_results) - success_count
            
            # 再次查询活动
            activity_after, _ = test.get_activity(activity_id)
            final_stock = activity_after.get('availableStock', 0)
            stock_reduction = orig_stock - final_stock
            
            print("\n=== 多用户并发购买测试结果 ===")
            print(f"总请求数: {len(all_results)}")
            print(f"成功请求数: {success_count}")
            print(f"失败请求数: {fail_count}")
            print(f"初始库存: {orig_stock}")
            print(f"最终库存: {final_stock}")
            print(f"库存减少量: {stock_reduction}")
            
            expected_success = min(registered, orig_stock)  # 应该最多有这么多成功
            
            if success_count <= expected_success:
                print("\n✅ 系统正确限制了每个用户只能购买一次")
            else:
                print(f"\n⚠️ 警告: 成功请求数({success_count})超过了预期({expected_success})!")
                print("系统可能存在重复购买风险!")
                
            if stock_reduction == success_count:
                print("\n✅ 库存减少量与成功订单数一致，系统工作正常")
            else:
                print(f"\n⚠️ 警告: 库存减少量({stock_reduction})与成功订单数({success_count})不一致!")
                
    except Exception as e:
        print(f"验证重复购买出错: {str(e)}")
        import traceback
        traceback.print_exc()

def verify_idempotency(test):
    try:
        activity_id = int(input("请输入活动ID: "))
        
        # 查询活动
        activity, _ = test.get_activity(activity_id)
        if not activity:
            print(f"未找到ID为 {activity_id} 的活动")
            return
            
        if activity.get('availableStock', 0) <= 0:
            print("活动库存已售罄，无法进行测试")
            return
            
        print(f"\n验证活动 '{activity.get('name')}' 的接口幂等性")
        
        # 注册测试用户
        print("注册测试用户...")
        prefix = f"idempotent_test_{int(time.time())}"
        registered = test.user_pool.register_batch(1, prefix)
        
        if registered <= 0:
            print("注册测试用户失败")
            return
        
        # 获取注册的用户ID - 从用户池中获取
        user_ids = list(test.user_pool.users.keys())
        if not user_ids:
            print("找不到已注册的用户")
            return
        
        # 使用最后注册的用户
        user_id = user_ids[-1]
        print(f"使用测试用户ID: {user_id}")
        
        # 第一次秒杀请求
        print("\n第一次秒杀请求...")
        success1, order_sn1, _ = test.seckill(user_id, activity_id)
        
        if not success1:
            print(f"首次秒杀失败，无法验证: {order_sn1}")
            return
            
        print(f"首次秒杀成功，订单号: {order_sn1}")
        
        # 查询库存
        activity_mid, _ = test.get_activity(activity_id)
        mid_stock = activity_mid.get('availableStock', 0)
        
        # 重复提交相同请求
        print(f"\n重复提交相同请求...")
        success2, order_sn2, _ = test.seckill(user_id, activity_id)
        
        # 再次查询库存
        activity_after, _ = test.get_activity(activity_id)
        after_stock = activity_after.get('availableStock', 0)
        
        print("\n幂等性测试结果:")
        print(f"首次请求后库存: {mid_stock}")
        print(f"重复请求后库存: {after_stock}")
        
        if mid_stock == after_stock:
            print("\n✅ 库存未变化，接口具有幂等性")
        else:
            print("\n⚠️ 警告: 重复请求导致库存再次减少，接口不具备幂等性!")
            
        if success2:
            if order_sn1 == order_sn2:
                print("✅ 重复请求返回相同订单号，接口具有幂等性")
            else:
                print(f"⚠️ 警告: 重复请求生成新订单 {order_sn2}，接口不具备幂等性!")
        else:
            print(f"重复请求被拒绝: {order_sn2}")
            
    except Exception as e:
        print(f"验证幂等性出错: {str(e)}")

def verify_idempotency_concurrent(test):
    try:
        activity_id = int(input("请输入活动ID: "))
        
        # 查询活动
        activity, _ = test.get_activity(activity_id)
        if not activity:
            print(f"未找到ID为 {activity_id} 的活动")
            return
            
        print(f"\n并发环境下验证活动 '{activity.get('name')}' 的接口幂等性")
        
        # 并发参数设置
        print("\n===== 并发幂等性测试选项 =====")
        print("1. 中等并发 (50用户/500并发)")
        print("2. 高并发 (100用户/1000并发)")
        print("3. 极限并发 (200用户/2000并发)")
        print("4. 自定义并发")
        print("5. 返回")
        
        option = input("\n选择并发级别: ")
        
        if option == "5":
            return
            
        if option == "1":
            user_count = 50
            concurrent_count = 500
        elif option == "2":
            user_count = 100
            concurrent_count = 1000
        elif option == "3":
            user_count = 200
            concurrent_count = 2000
        elif option == "4":
            try:
                user_count = int(input("测试用户数: "))
                concurrent_count = int(input("总并发请求数: "))
            except ValueError:
                print("输入无效，使用默认值")
                user_count = 50
                concurrent_count = 500
        else:
            print("选择无效，使用中等并发设置")
            user_count = 50
            concurrent_count = 500
        
        print(f"\n将执行并发幂等性测试:")
        print(f"- 测试用户数: {user_count}")
        print(f"- 总并发请求数: {concurrent_count}")
        print(f"- 每用户平均并发请求: {concurrent_count//user_count}")
            
        confirm = input("\n确认执行测试? (y/n): ")
        if confirm.lower() != 'y':
            print("已取消测试")
            return
            
        # 注册测试用户
        print(f"注册{user_count}个测试用户...")
        registered = test.user_pool.register_batch(user_count, f"idem_test_{int(time.time())}")
        print(f"成功注册{registered}个用户")
        
        if registered < 10:
            print("注册用户数太少，无法进行有效测试")
            return
        
        # 获取用户ID列表
        user_ids = list(test.user_pool.users.keys())[-registered:]
        
        # 记录初始库存
        init_stock = activity.get('availableStock', 0)
        print(f"初始库存: {init_stock}")
        
        # 创建并发测试函数
        def send_request(user_id, activity_id):
            return test.seckill(user_id, activity_id)
        
        # 准备所有请求
        all_tasks = []
        for _ in range(concurrent_count):
            # 随机选择用户，每个用户会被多次选中，模拟重复提交
            user_id = random.choice(user_ids)
            all_tasks.append((user_id, activity_id))
        
        print(f"\n开始执行{concurrent_count}个并发请求...")
        start_time = time.time()
        
        # 使用线程池执行所有请求
        results = []
        with concurrent.futures.ThreadPoolExecutor(max_workers=min(concurrent_count, 500)) as executor:
            futures = [executor.submit(send_request, user_id, act_id) for user_id, act_id in all_tasks]
            
            for future in concurrent.futures.as_completed(futures):
                try:
                    result = future.result()
                    results.append(result)
                except Exception as e:
                    print(f"请求执行错误: {str(e)}")
        
        end_time = time.time()
        
        # 计算每个用户的成功率
        user_success = {}  # user_id -> [订单号]
        for i, (success, msg, _) in enumerate(results):
            user_id = all_tasks[i][0]
            if success:
                if user_id not in user_success:
                    user_success[user_id] = []
                user_success[user_id].append(msg)  # msg是订单号
        
        # 分析结果
        success_count = sum(1 for r in results if r[0])
        fail_count = len(results) - success_count
        
        # 查询最终库存
        activity_after, _ = test.get_activity(activity_id)
        final_stock = activity_after.get('availableStock', 0)
        stock_reduction = init_stock - final_stock
        
        print("\n===== 并发幂等性测试结果 =====")
        print(f"测试用时: {end_time - start_time:.2f}秒")
        print(f"总请求数: {len(results)}")
        print(f"成功请求数: {success_count}")
        print(f"失败请求数: {fail_count}")
        print(f"QPS: {len(results)/(end_time - start_time):.2f}次/秒")
        print(f"\n初始库存: {init_stock}")
        print(f"最终库存: {final_stock}")
        print(f"库存减少量: {stock_reduction}")
        print(f"成功下单用户数: {len(user_success)}")
        
        # 检查每个用户是否只有一个订单
        multi_orders = 0
        for user_id, orders in user_success.items():
            if len(set(orders)) > 1:
                multi_orders += 1
                print(f"⚠️ 用户{user_id}有多个不同订单: {set(orders)}")
        
        if multi_orders == 0:
            print("\n✅ 所有用户都只有一个订单，幂等性检查通过")
        else:
            print(f"\n⚠️ 警告: 有{multi_orders}个用户存在多个订单!")
            
        if stock_reduction == len(user_success):
            print("\n✅ 库存减少量与成功下单用户数一致，接口具有幂等性")
        else:
            print(f"\n⚠️ 警告: 库存减少量({stock_reduction})与成功下单用户数({len(user_success)})不一致!")
            print("接口在高并发下可能存在幂等性问题!")
            
    except Exception as e:
        print(f"并发幂等性测试出错: {str(e)}")
        import traceback
        traceback.print_exc()

def change_server(test):
    current_url = test.base_url
    print(f"当前服务器: {current_url}")
    
    new_url = input(f"输入新服务器地址 [{current_url}]: ") or current_url
    
    if new_url != current_url:
        test.base_url = new_url
        test.session.headers.update({'Content-Type': 'application/json'})
        test.user_pool = UserPool(new_url)
        print(f"服务器地址已更新为: {new_url}")
    else:
        print("服务器地址未变更")

if __name__ == "__main__":
    try:
        simple_menu()
    except KeyboardInterrupt:
        print("\n程序已退出")
    except Exception as e:
        print(f"程序执行异常: {str(e)}")