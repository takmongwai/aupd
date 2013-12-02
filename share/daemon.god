# encoding: UTF-8

app_name = "aupd_daemon"
app_file = "aupd"
app_path = "/srv/aupd"
log_file = "#{app_path}/log/#{app_name}.log"
pid_file = "#{app_path}/tmp/#{app_name}.pid"

#运行进程的用户和组
app_user = "www-data"
#启动命令
start_command = "\
cd #{app_path} && \
#{app_path}/bin/#{app_file}"

#停止命令
stop_command = "kill -QUIT `cat #{pid_file}`"
#重启命令
restart_command = "#{stop_command} && #{start_command}"

God.watch do |w|
  #分组名
  #w.group = "#{app_name}"
  w.dir = app_path
  w.name = "#{app_name}"
  w.env = {"GOMAXPROCS" => 8}
  w.log = log_file
  w.pid_file = pid_file
  w.interval = 30.seconds
  w.start = start_command
  w.stop = stop_command
  w.restart = restart_command
  #缓冲时间
  w.start_grace = 10.seconds
  w.restart_grace = 10.seconds
  w.behavior(:clean_pid_file)

  w.start_if do |start|
    start.condition(:process_running) do |c|
      c.interval = 5.seconds
      c.running = false
    end
  end

  w.restart_if do |restart|
    restart.condition(:memory_usage) do |c|
      c.above = 500.megabytes
      c.times = [3, 5]
    end
    restart.condition(:cpu_usage) do |c|
      c.above = 90.percent
      c.times = 20
    end
  end
  w.lifecycle do |on|
    on.condition(:flapping) do |c|
      c.to_state = [:start, :restart]
      c.times = 5
      c.within = 5.minute
      c.transition = :unmonitored
      c.retry_in = 10.minutes
      c.retry_times = 5
      c.retry_within = 2.hours
    end
  end
end
