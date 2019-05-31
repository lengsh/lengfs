{{ template "header" . }}
<div>
<h1>Hello,lengfs !</h1>
</div>

Parent = {{.Node.Parent}}<BR>
Inode = {{.Node.Inode}}<BR>
Pnode = {{.Node.Pnode}}<BR>
Start = {{DateFormat .Stat.StartTime "2006-01-02 15:04:05"}}<BR>
Modetime = {{DateFormat .Stat.ModTime "2006-01-02 15:04:05"}}<BR>
<h2>Command</h2>
<div class="divcss5">
<a href="/lfs/"></a>
<BR><a href="/lfs/uptest/">upfile</a>
<BR><a href="/lfs/upload/">up result</a>
<BR><a href="/lfs/pathinfo/?date={{DateFormat .Stat.ModTime "20060102" }}&inode={{.Node.Inode}}">pathinfo </a>
<BR><a href="/lfs/psync/?">psync</a>
</div>
<BR>
{{ template "footer" . }}
