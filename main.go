package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>HTTP 重定向演示</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            line-height: 1.6;
        }
        .redirect-type {
            margin: 20px 0;
            padding: 15px;
            border: 1px solid #ddd;
            border-radius: 5px;
        }
        h2 {
            color: #333;
        }
        .description {
            background-color: #f5f5f5;
            padding: 10px;
            border-radius: 3px;
        }
        .test-link {
            display: inline-block;
            margin: 10px 0;
            padding: 8px 15px;
            background-color: #007bff;
            color: white;
            text-decoration: none;
            border-radius: 4px;
        }
        .test-link:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <h1>HTTP 重定向状态码演示</h1>
    
    <div class="redirect-type">
        <h2>301 - Moved Permanently（永久重定向）</h2>
        <div class="description">
            <p>用途：资源永久移动到新位置</p>
            <p>特点：</p>
            <ul>
                <li>浏览器会缓存这个重定向</li>
                <li>搜索引擎会更新资源的新地址</li>
                <li>通常用于网站域名变更或者永久性的URL结构改变</li>
            </ul>
        </div>
        <a href="/redirect/301" class="test-link">测试 301 重定向</a>
    </div>

    <div class="redirect-type">
        <h2>302 - Found（临时重定向）</h2>
        <div class="description">
            <p>用途：资源临时移动到新位置</p>
            <p>特点：</p>
            <ul>
                <li>浏览器不会缓存重定向</li>
                <li>保持原有的 HTTP 方法</li>
                <li>常用于临时维护页面、A/B测试等场景</li>
                <li>在用户登录场景中广泛使用：
                    <ul>
                        <li>未登录用户访问需要认证的页面时，重定向到登录页面</li>
                        <li>登录成功后，重定向回原始请求页面</li>
                    </ul>
                </li>
            </ul>
        </div>
        <div style="margin-top: 15px;">
            <a href="/protected-page" class="test-link">访问需要登录的页面（302 重定向到登录）</a>
        </div>
    </div>

    <div class="redirect-type">
        <h2>303 - See Other</h2>
        <div class="description">
            <p>用途：通常用于POST请求后的重定向</p>
            <p>特点：</p>
            <ul>
                <li>强制使用GET方法访问新地址</li>
                <li>常用于表单提交后重定向到结果页面</li>
                <li>防止表单重复提交</li>
            </ul>
        </div>
        <form action="/submit-form" method="POST" style="margin-top: 15px;">
            <input type="text" name="message" placeholder="输入一些内容..." style="padding: 8px; margin-right: 10px; border: 1px solid #ddd; border-radius: 4px;">
            <button type="submit" class="test-link" style="border: none; cursor: pointer;">提交表单 (303 重定向)</button>
        </form>
    </div>

    <div class="redirect-type">
        <h2>307 - Temporary Redirect</h2>
        <div class="description">
            <p>用途：临时重定向，但严格保持原有HTTP方法</p>
            <p>特点：</p>
            <ul>
                <li>与302类似，但保证请求方法不变</li>
                <li>如果原请求是POST，重定向后仍然是POST</li>
                <li>适用于需要严格保持HTTP方法的场景</li>
            </ul>
        </div>
        <form action="/redirect/307" method="POST" style="margin-top: 15px;">
            <input type="text" name="message" placeholder="输入一些内容..." style="padding: 8px; margin-right: 10px; border: 1px solid #ddd; border-radius: 4px;">
            <button type="submit" class="test-link" style="border: none; cursor: pointer;">POST 提交 (307 重定向)</button>
        </form>
    </div>

    <div class="redirect-type">
        <h2>308 - Permanent Redirect</h2>
        <div class="description">
            <p>用途：永久重定向，但严格保持原有HTTP方法</p>
            <p>特点：</p>
            <ul>
                <li>类似于301，但保证请求方法不变</li>
                <li>永久性质的重定向</li>
                <li>适用于需要永久重定向且保持HTTP方法的场景</li>
            </ul>
        </div>
        <div style="margin-top: 15px;">
            <input type="text" id="put-message" placeholder="输入一些内容..." style="padding: 8px; margin-right: 10px; border: 1px solid #ddd; border-radius: 4px;">
            <button onclick="sendPutRequest()" class="test-link" style="border: none; cursor: pointer;">PUT 提交 (308 重定向)</button>
        </div>
        <script>
        function sendPutRequest() {
            const message = document.getElementById('put-message').value;
            fetch('/redirect/308', {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: 'message=' + encodeURIComponent(message),
                redirect: 'follow'
            }).then(response => {
                if (response.redirected) {
                    window.location.href = response.url;
                }
            });
        }
        </script>
    </div>
</body>
</html>
`

func main() {
	// 创建模板
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		log.Fatal(err)
	}

	// 首页处理
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		tmpl.Execute(w, nil)
	})

	// 目标页面
	http.HandleFunc("/destination", func(w http.ResponseWriter, r *http.Request) {
		redirectType := r.URL.Query().Get("type")
		message := r.URL.Query().Get("message")
		messageHTML := ""
		if message != "" {
			messageHTML = fmt.Sprintf(`<p>提交的消息: %s</p>`, message)
		}

		fmt.Fprintf(w, `
			<html>
			<body>
				<h1>目标页面</h1>
				<p>你是通过 %s 重定向到达此页面的</p>
				<p>使用的 HTTP 方法: %s</p>
				%s
				<p><a href="/">返回首页</a></p>
			</body>
			</html>
		`, redirectType, r.Method, messageHTML)
	})

	// 表单提交处理
	http.HandleFunc("/submit-form", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "只接受 POST 请求", http.StatusMethodNotAllowed)
			return
		}

		message := r.FormValue("message")
		// 将消息存储在查询参数中，这样重定向后的页面可以显示它
		destination := fmt.Sprintf("/destination?type=303&message=%s", message)
		http.Redirect(w, r, destination, http.StatusSeeOther)
	})

	// 检查用户登录状态
	checkLogin := func(r *http.Request) (loggedIn bool, username string) {
		cookie, err := r.Cookie("session")
		if err != nil {
			return false, ""
		}
		return true, cookie.Value
	}

	// 受保护的页面
	http.HandleFunc("/protected-page", func(w http.ResponseWriter, r *http.Request) {
		loggedIn, username := checkLogin(r)
		if !loggedIn {
			// 保存用户想要访问的原始 URL
			returnTo := r.URL.String()
			// 重定向到登录页面，带上返回地址
			http.Redirect(w, r, fmt.Sprintf("/login?return_to=%s", returnTo), http.StatusFound) // 302 重定向
			return
		}
		// 已登录用户可以看到受保护的内容
		fmt.Fprintf(w, `
			<html>
			<body>
				<h1>受保护的页面</h1>
				<p>欢迎，%s！</p>
				<p>这是一个需要登录才能访问的页面。</p>
				<p>你现在已经登录了，所以可以看到这个内容。</p>
				<p><a href="/logout">退出登录</a></p>
				<p><a href="/">返回首页</a></p>
			</body>
			</html>
		`, username)
	})

	// 登录页面
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// 获取用户名
			username := r.FormValue("username")
			if username == "" {
				http.Error(w, "请输入用户名", http.StatusBadRequest)
				return
			}

			// 设置登录 cookie
			cookie := &http.Cookie{
				Name:     "session",
				Value:    username,
				Path:     "/",
				MaxAge:   3600, // 1小时过期
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			
			// 获取登录成功后要返回的页面
			returnTo := r.FormValue("return_to")
			if returnTo == "" {
				returnTo = "/protected-page"
			}
			
			// 登录成功后重定向到原始请求的页面
			http.Redirect(w, r, returnTo, http.StatusFound) // 302 重定向
			return
		}

		// 显示登录表单
		returnTo := r.URL.Query().Get("return_to")
		fmt.Fprintf(w, `
			<html>
			<body>
				<h1>登录页面</h1>
				<p>这是一个模拟的登录页面，展示 302 重定向在登录流程中的应用。</p>
				<form method="POST" action="/login" style="margin: 20px 0;">
					<input type="hidden" name="return_to" value="%s">
					<div style="margin-bottom: 10px;">
						<input type="text" name="username" placeholder="输入用户名" 
							style="padding: 8px; border: 1px solid #ddd; border-radius: 4px; width: 200px;">
					</div>
					<button type="submit" style="padding: 8px 15px; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer;">
						登录
					</button>
				</form>
				<p><a href="/">返回首页</a></p>
			</body>
			</html>
		`, returnTo)
	})

	// 退出登录
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		// 删除 cookie
		cookie := &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // 立即删除 cookie
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/", http.StatusFound) // 302 重定向
	})

	// 重定向处理器
	http.HandleFunc("/redirect/", func(w http.ResponseWriter, r *http.Request) {
		redirectType := r.URL.Path[len("/redirect/"):]
		
		// 处理 POST 和 PUT 请求数据
		message := ""
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			r.ParseForm()
			message = r.FormValue("message")
		}

		// 构建目标 URL
		destination := fmt.Sprintf("/destination?type=%s", redirectType)
		if message != "" {
			destination = fmt.Sprintf("%s&message=%s", destination, message)
		}

		switch redirectType {
		case "301":
			http.Redirect(w, r, destination, http.StatusMovedPermanently)
		case "302":
			http.Redirect(w, r, destination, http.StatusFound)
		case "307":
			http.Redirect(w, r, destination, http.StatusTemporaryRedirect)
		case "308":
			http.Redirect(w, r, destination, http.StatusPermanentRedirect)
		default:
			http.NotFound(w, r)
		}
	})

	// 启动服务器
	fmt.Println("服务器启动在 http://localhost:8060")
	log.Fatal(http.ListenAndServe(":8060", nil))
}
