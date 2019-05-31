{{ template "header" . }}

 <div data-role="main" class="ui-content">
  <div class="ui-field-contain">
  <form action="/lfs/upload/" method="post" enctype="multipart/form-data" data-ajax="false">
  <label for="upfile">图片:</label>
  <input  id="uploadimg" name="file"  type="file"  runat="server" method="post" enctype="multipart/form-data" data-inline="true"  data-ajax="false" /> 
 <label for="desc">说明:</label>
 <textarea name="desc" id="desc"></textarea>
     <input type="submit" value="upload" >
    <input type="hidden" name="goto" value="/lfs/uptest/"  />
     <input type="hidden" name=".scrumb" value="{{ Scrumb }}"  >
 </form>
 </div>
 </div>

<div align=center>
{{ range . }}
<img src="{{.Thumbnail}}" /><BR>
{{.Original}}

{{ end }}
</div>

{{ template "footer" . }}
