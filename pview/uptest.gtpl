{{ template "header" . }}

 <div align=center>
<BR>
<form enctype="multipart/form-data" action="/lfs/upload/" method="post" >
      <input type="file" name="file" />
      <input type="submit" value="upload" />
      <!-- input type="hidden" name="goto" value="/lfs/uptest/"  / -->
      <input type="hidden" name=".scrumb" value="{{Scrumb}}"  />
 </form>
 </div>
<BR>
<div align=center>
{{ range . }}
<img src="{{.Thumbnail}}" /><BR>
{{.Original}}<BR>
{{ end }}
</div>

{{ template "footer" . }}
