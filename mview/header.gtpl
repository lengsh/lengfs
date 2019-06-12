{{define "header"}}

<!DOCTYPE html>
<html>
<head>
<meta name="viewport" content="width=device-width, initial-scale=1"  >
<link rel="stylesheet" href="https://apps.bdimg.com/libs/jquerymobile/1.4.5/jquery.mobile-1.4.5.min.css">
<script src="https://apps.bdimg.com/libs/jquery/1.10.2/jquery.min.js"></script>
<script src="https://apps.bdimg.com/libs/jquerymobile/1.4.5/jquery.mobile-1.4.5.min.js"></script>
<style>
.divcss5{padding-left:50px;margin-top:-22px}
</style>
</head>
<body>
<div data-role="page" id="pageone">
   <div data-role="header">
      <a href="/lfs/" class="ui-btn ui-icon-home ui-btn-icon-left">主页</a>
      <h1>欢迎{{UserName}}</h1>
      <a href="/lfs/login" class="ui-btn ui-icon-search ui-btn-icon-left">登录</a>
   </div>

{{end}}
