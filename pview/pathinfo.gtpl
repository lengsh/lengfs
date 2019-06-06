{{ template "header" . }}
<div align=center>
<h1>Hello,lengfs !</h1>
</div>
<div style="position:relative;left:100px;top:10px">
{{$root:=.root}}
{{range .data}}
{{if IsContains . ".thumbnail."}}
<img src="/{{$root}}{{.}}" /><BR>
{{else}}
{{.}}<BR>
{{end}}
{{end}}
<BR>
</div>
{{ template "footer" . }}
