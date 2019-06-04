{{ template "header" . }}
<div>
<h1>
{{if .user}}
Hello,{{.user}}
{{else}}
Hello,anonymouse!
{{end}}
</h1>
</div>

<div data-role="main" class="ui-content">
 <form method="post" action="/lfs/login">
        <div>
          <h3>登录信息</h3>
          <label for="usrnm" class="ui-hidden-accessible">用户名:</label>
          <input type="text" name="email" id="usrnm" placeholder="用户名">
          <label for="pswd" class="ui-hidden-accessible">密码:</label>
          <input type="password" name="password" id="pswd" placeholder="密码">
          <input type="submit" data-inline="true" value="登录">
        </div>
  </form>
</div>


{{ template "footer" . }}
