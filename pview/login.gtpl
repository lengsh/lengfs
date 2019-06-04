{{ template "header" . }}
<div>
    <H1>Welcome 
    {{if .user}} 
    {{.user}}
    {{else}}
    Anonymouse
   {{end}}
   </h1>
</div>

<body class="grey lighten-4">
<div class="admin-ui row">
<div class="init col s5">
<div class="card">
<div class="card-content">
    <blockquote>请输入登录账号和密码.</blockquote>
    <form method="post" action="/lfs/login" class="row">
        <div class="input-field col s12">
            <input placeholder="Enter your email address e.g. you@example.com" class="validate required" type="email" id="email" name="email"/>
            <label for="email" class="active">Email</label>
        </div>
        <div class="input-field col s12">
            <input placeholder="Enter your password" class="validate required" type="password" id="password" name="password"/>
            <a href="/recover">Forgot password?</a>            
            <label for="password" class="active">Password</label>  
        </div>
        <button class="btn waves-effect waves-light right">Log in</button>
    </form>
</div>
</div>
</div>
</div>

{{ template "footer" . }}
