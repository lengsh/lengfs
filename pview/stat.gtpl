{{ template "header" . }}
<div align=center>
<h1>Hello,{{UserName}}!</h1>
</div>
<div style="position:relative;left:100px;top:10px">
Parent = {{.Node.Parent}}<BR>
Inode = {{.Node.Inode}}<BR>
Pnode = {{.Node.Pnode}}<BR>
Server Queues = {{.Node.Queues}}<BR>
Start = {{DateFormat .Stat.StartTime "2006-01-02 15:04:05"}}<BR>
Modetime = {{DateFormat .Stat.ModTime "2006-01-02 15:04:05"}}<BR>
Disk Used = <font color="red">{{.Stat.Used}}</font><BR>
Disk Scan Frequency = {{.Stat.Frequency}}秒<BR>
<h2>Command</h2>
<div>
<a href="/lfs/"></a>
<BR><a href="/lfs/uptest/">upfile</a>
<BR><a href="/lfs/upload/">up result</a>
<BR><a href="/lfs/pathinfo/?date={{DateFormat .Stat.ModTime "20060102" }}&inode={{.Node.Inode}}">pathinfo </a>
<BR><a href="/lfs/pathinfo/?date={{DateFormat .Stat.ModTime "20060102" }}&inode={{.Node.Inode}}&&pathinfo=Yes">pathinfo(Meta) </a>
<BR><a href="/lengfs/">lengfs</a>
</div>
<BR>
</div>
{{ template "footer" . }}
