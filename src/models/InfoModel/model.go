package InfoModel

import (
	"strings"
)

var notes map[string]string

func init() {
	notes = make(map[string]string, 0)

	//Server
	notes["# Server"] = "服务端信息"
	notes["redis_version"] = "Redis服务器版本"
	notes["redis_git_sha1"] = "Redis在Git上的SHA1值"
	notes["redis_git_dirty"] = "Git dirty标志"
	notes["redis_build_id"] = "Redis构建ID"
	notes["redis_mode"] = "运行模式，单机或者集群"
	notes["os"] = "Redis服务器的宿主操作系统"
	notes["arch_bits"] = "架构（32 或 64 位）"
	notes["multiplexing_api"] = "Redis所使用的事件处理机制"
	notes["atomicvar_api"] = "原子变量API"
	notes["gcc_version"] = "编译redis时所使用的gcc版本"
	notes["process_id"] = "Redis服务器进程的pid"
	notes["process_supervised"] = "监督进程"
	notes["run_id"] = "Redis服务器的随机标识符(用于sentinel和集群)"
	notes["tcp_port"] = "Redis服务器监听端口"
	notes["server_time_usec"] = "最后一次访问的时间戳，单位是微秒"
	notes["uptime_in_seconds"] = "Redis服务器启动总时间，单位是秒"
	notes["uptime_in_days"] = "Redis服务器启动总时间，单位是天"
	notes["hz"] = "Redis内部调度（进行关闭timeout的客户端，删除过期key等等）频率，程序规定serverCron每秒运行10次。"
	notes["configured_hz"] = "已经配置的hz，如上"
	notes["lru_clock"] = "位置自增的时钟，用于LRU管理,该时钟100ms(hz=10,因此每1000ms/10=100ms执行一次定时任务)更新一次。"
	notes["executable"] = "执行文件"
	notes["config_file"] = "配置文件路径"
	notes["io_threads_active"] = "IO线程是否活跃"

	//Clients
	notes["# Clients"] = "客户端信息"
	notes["connected_clients"] = "已经连接的客户端数"
	notes["cluster_connections"] = "集群连接数"
	notes["maxclients"] = "最大连接客户端数"
	notes["client_recent_max_input_buffer"] = "当前连接客户端中最长的输出列表"
	notes["client_recent_max_output_buffer"] = "当前连接的客户端当中，最大输入缓存"
	notes["blocked_clients"] = "正在等待阻塞命令(BLPOP、BRPOP、BRPOPLPUSH)的客户端的数量"
	notes["tracking_clients"] = "正在跟踪的客户端"
	notes["clients_in_timeout_table"] = "在超时表中的客户端"

	//Memory
	notes["# Memory"] = "内存"
	notes["used_memory"] = "已使用内存，以字节（byte）为单位以"
	notes["used_memory_human"] = "人类可读的格式返回Redis分配的内存总量"
	notes["used_memory_rss"] = "系统给redis分配的内存即常驻内存，和top、ps 等命令的输出一致。"
	notes["used_memory_rss_human"] = "以人类可读的格式返回常驻内存"
	notes["used_memory_peak"] = "内存使用的峰值大小"
	notes["used_memory_peak_human"] = "以人类可读的格式返回内存使用的峰值大小"
	notes["used_memory_peak_perc"] = "使用内存达到峰值内存的百分比"
	notes["used_memory_overhead"] = "Redis为了维护数据集的内部机制所需的内存开销，包括所有客户端输出缓冲区、查询缓冲区、AOF重写缓冲区和主从复制的backlog。"
	notes["used_memory_startup"] = "Redis启动完成使用的内存"
	notes["used_memory_dataset"] = "数据占用内存大小 等于 used_memory-user_memory_overhead"
	notes["used_memory_dataset_perc"] = "数据占用的内存大小百分比,(used_memory_dataset / (used_memory - used_memory_startup))*100%"
	notes["allocator_allocated"] = "分配器分配的内存"
	notes["allocator_active"] = "分配器活跃的内存"
	notes["allocator_resident"] = "分配器常驻的内存"
	notes["total_system_memory"] = "系统总内存"
	notes["total_system_memory_human"] = "人类方式直观展示系统总内存"
	notes["used_memory_lua"] = "lua脚本存储占用的内存"
	notes["used_memory_lua_human"] = "人类方式直观展示lua脚本存储占用的内存"
	notes["used_memory_scripts"] = "脚本存储占用的内存"
	notes["used_memory_scripts_human"] = "人类方式直观展示脚本存储占用的内存"
	notes["number_of_cached_scripts"] = "被缓存的脚本总数"
	notes["maxmemory"] = "配置允许redis使用的最大内存，0为不限制"
	notes["maxmemory_human"] = "人类方式显示允许redis使用的最大内存，0为不限制"
	notes["maxmemory_policy"] = "配置的内存淘汰策略,当前配置为内存满后写入失败"
	notes["allocator_frag_ratio"] = "分配器的碎片率"
	notes["allocator_frag_bytes"] = "分配器的碎片大小"
	notes["allocator_rss_ratio"] = "分配器常驻内存比例"
	notes["allocator_rss_bytes"] = "分配器的常驻内存大小"
	notes["rss_overhead_ratio"] = "常驻内存开销比例"
	notes["rss_overhead_bytes"] = "常驻内存开销大小"
	notes["mem_fragmentation_ratio"] = "碎片率(used_memory_rss / used_memory),正常(1,1.6),大于比例说明内存碎片严重"
	notes["mem_fragmentation_bytes"] = "内存碎片大小"
	notes["mem_not_counted_for_evict"] = "被驱逐的内存"
	notes["mem_replication_backlog"] = "Redis复制积压缓冲区内存"
	notes["mem_clients_slaves"] = "Redis节点客户端消耗内存"
	notes["mem_clients_normal"] = "Redis所有常规客户端消耗内存"
	notes["mem_aof_buffer"] = "AOF使用内存"
	notes["mem_allocator"] = "内存分配器"
	notes["active_defrag_running"] = "是否有磁盘碎片正在运行"
	notes["lazyfree_pending_objects"] = "活动碎片整理是否处于活动状态(0没有,1正在运行)"
	notes["lazyfreed_objects"] = "0-不存在延迟释放的挂起对象"

	//Persistence
	notes["# Persistence"] = "持久化信息"
	notes["loading"] = "表示是否正在载入持久化文件"
	notes["rdb_changes_since_last_save"] = "距离最近一次成功创建RDB持久化文件之后，经过了多少秒"
	notes["rdb_bgsave_in_progress"] = "表示是否正在执行bgsave创建持久化文件"
	notes["rdb_last_save_time"] = "最近一次成功创建RDB文件的时间戳"
	notes["rdb_last_bgsave_status"] = "表示最近一次创建RDB文件是否成功"
	notes["rdb_last_bgsave_time_sec"] = "记录最近一次创建RDB文件耗费的秒数"
	notes["rdb_current_bgsave_time_sec"] = "当前正在创建RDB文件已经耗费的秒数"
	notes["rdb_last_cow_size"] = "记录父进程与子进程相比执行了多少修改"
	notes["aof_enabled"] = "是否开启AOF"
	notes["aof_rewrite_in_progress"] = "记录是否正在创建AOF文件"
	notes["aof_rewrite_scheduled"] = "记录在AOF文件创建完毕后，是否需要执行预约的AOF重写操作"
	notes["aof_last_rewrite_time_sec"] = "最近一次创建AOF文件耗费的时长"
	notes["aof_current_rewrite_time_sec"] = "当前正在创建AOF文件已经耗费的秒数"
	notes["aof_last_bgrewrite_status"] = "表示最近一次子进程创建AOF文件是否成功"
	notes["aof_last_write_status"] = "表示最近一次主进程创建AOF文件是否成功"
	notes["aof_last_cow_size"] = "记录父进程与子进程相比执行了多少修改"
	notes["module_fork_in_progress"] = "模块分叉在进程中"
	notes["module_fork_last_cow_size"]="在最后一个模块fork操作期间，以字节为单位的写时拷贝内存的大小"
	notes["aof_current_size"] = "AOF文件当前大小"
	notes["aof_base_size"] = "服务启动时或AOF重写最近一次执行后AOF文件大小"
	notes["aof_pending_rewrite"] = "是否有AOF重写操作在等待执行"
	notes["aof_buffer_length"] = "AOF缓冲区大小"
	notes["aof_rewrite_buffer_length"] = "AOF重写缓冲区大小"
	notes["aof_pending_bio_fsync"] = "一种Redis AOF刷新策略"
	notes["aof_delayed_fsync"] = "原主数据库追加aof阻塞"

	//# Stats

	notes["# Stats"] = "通用统计数据"
	notes["total_connections_received"] = "服务器已接收连接的数量"
	notes["total_commands_processed"] = "服务器已执行命令数量"
	notes["instantaneous_ops_per_sec"] = "服务器每秒执行命令数量OPS"
	notes["total_net_input_bytes"] = "网络输入总字节数"
	notes["total_net_output_bytes"] = "网络输出总字节数"
	notes["instantaneous_input_kbps"] = "网络输入kbps"
	notes["instantaneous_output_kbps"] = "网络输出kbps"
	notes["rejected_connections"] = "因最大客户端数量限制而被拒绝的连接请求数量"
	notes["sync_full"] = "主从全量同步成功次数"
	notes["sync_partial_ok"] = "主从部分同步成功次数"
	notes["sync_partial_err"] = "主从部分同步失败次数"
	notes["expired_keys"] = "运行以来过期的key数量"
	notes["expired_stale_perc"] = "过期的比率"
	notes["expired_time_cap_reached_count"] = "因cpu利用率过高而提前结束定期删除的次数"
	notes["expire_cycle_cpu_milliseconds"] = "在活动到期周期上花费的累计时间量"
	notes["evicted_keys"] = "运行以来剔除(超过maxmemory后)的key的数量"
	notes["keyspace_hits"] = "在主字典中成功查找键的数量"
	notes["keyspace_misses"] = "主词典中键查找失败的次数"
	notes["pubsub_channels"] = "当前使用中的频道的数量"
	notes["pubsub_patterns"] = "当前使用中的模式的数量"
	notes["latest_fork_usec"] = "最近一次fork操作阻塞redis进程的耗时，单位：微秒"
	notes["total_forks"] = "总共的forks操作"
	notes["migrate_cached_sockets"] = "是否已缓存了到改地址的连接"
	notes["slave_expires_tracked_keys"] = "主从实例到期key数量"
	notes["active_defrag_hits"] = "主动碎片整理命中次数"
	notes["active_defrag_misses"] = "主动碎片整理未命中次数"
	notes["active_defrag_key_hits"] = "主动碎片整理key命中次数"
	notes["active_defrag_key_misses"] = "主动碎片整理key为命中次数"
	notes["tracking_total_keys"] = "所有跟踪的键盘"
	notes["tracking_total_items"] = "所有跟踪的项目"
	notes["tracking_total_prefixes"] = "追踪总计的前缀"
	notes["unexpected_error_replies"] = "意外的错误回复"
	notes["dump_payload_sanitizations"] = "导出负载清理"
	notes["total_reads_processed"] = "处理的读取事件总数"
	notes["total_writes_processed"] = "处理的写入事件总数"
	notes["io_threaded_reads_processed"] = "主线程和I / O线程处理的读取事件数"
	notes["io_threaded_writes_processed"] = "主线程和I / O线程处理的写事件数"

	//# Replication
	notes["# Replication"] = "异步复制"
	notes["role"] = "角色"
	notes["connected_slaves"] = "连接的Slave的数量。"
	notes["master_replid"] = "master启动时生成的40位16进制的随机字符串，用来标识master节点"
	notes["master_replid2"] = "master启动时生成的40位16进制的随机字符串，用来标识master节点2"
	notes["master_repl_offset"] = "复制流中的一个偏移量，master处理完写入命令后，会把命令的字节长度做累加记录，统计在该字段。该字段也是实现部分复制的关键字段。"
	notes["second_repl_offset"] = "同样也是一个偏移量，从节点收到主节点发送的命令后，累加自身的偏移量，通过比较主从节点的复制偏移量可以判断主从节点数据是否一致。"
	notes["repl_backlog_active"] = "是否开启了backlog。"
	notes["repl_backlog_size"] = "复制积压缓冲区的总大小（字节）"
	notes["repl_backlog_first_byte_offset"] = "backlog中保存的Master最早的偏移量"
	notes["repl_backlog_histlen"] = "复制积压缓冲区中数据的大小（字节）"

	//# CPU
	notes["# CPU"] = "CPU信息"
	notes["used_cpu_sys"] = "被redis服务端消耗的系统CPU"
	notes["used_cpu_user"] = "被redis服务端消耗的用户CPU"
	notes["used_cpu_sys_children"] = "被后台程序消耗的系统CPU"
	notes["used_cpu_user_children"] = "被后台程序消耗的用户CPU"
	notes["used_cpu_sys_main_thread"] = "主线程种被redis服务端消耗的系统CPU"
	notes["used_cpu_user_main_thread"] = "主线程种被redis服务端消耗的用户CPU"

	//# Modules
	notes["# Modules"] = "模块"

	//# Cluster
	notes["# Cluster"] = "集群"
	notes["cluster_enabled"] = "集群是否开启"

	//# Keyspace
	notes["# Keyspace"] = "关键的空间"
	notes["db0"] = "数据库0"

}

type Info struct {
	Key   string `json:"key"`
	Value string `json:"value"`

	Desc string `json:"desc"`
}

func NewInfo() *Info {
	return &Info{}
}

func GetInfo(info string) []*Info {

	infos := make([]*Info, 0)
	split := strings.Split(info, "\r\n")

	for i := 0; i < len(split); i++ {
		v := split[i]
		info := NewInfo()
		if strings.Index(v, ":") != -1 {
			k := strings.Split(v, ":")[0]
			v := strings.Split(v, ":")[1]
			info.Key = k
			info.Value = v
			info.Desc = notes[k]
		} else {
			info.Key = v
			info.Value = ""
			if strings.TrimSpace(v) != "" {
				v = strings.Replace(v, string(0x0d), "", -1)
				v = strings.Replace(v, string(0x0a), "", -1)
				info.Desc = notes[v]
			} else {
				info.Desc = ""
			}
		}
		infos = append(infos, info)
	}

	return infos
}
