package main

import (
	"context"
	"log"
	"net/http"
	"strings"
)

// SubdomainContext はサブドメイン情報を保持する構造体
type SubdomainContext struct {
	Subdomain string
	Domain    string
}

// ExtractSubdomain はリクエストからサブドメインを抽出する
func ExtractSubdomain(r *http.Request) *SubdomainContext {
	host := r.Host

	// ポート番号を除去
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}

	// サブドメインを抽出
	parts := strings.Split(host, ".")

	if len(parts) >= 3 && strings.HasSuffix(host, "fumi042-server.top") {
		return &SubdomainContext{
			Subdomain: parts[0],
			Domain:    host,
		}
	}

	return nil
}

// SubdomainMiddleware はサブドメインに基づいてリクエストを処理するミドルウェア
func SubdomainMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subdomainCtx := ExtractSubdomain(r)

		if subdomainCtx != nil {
			log.Printf("Subdomain detected: %s", subdomainCtx.Subdomain)

			// サブドメイン情報をコンテキストに追加
			ctx := r.Context()
			ctx = context.WithValue(ctx, "subdomain", subdomainCtx)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// SubdomainRouter はサブドメインに基づいてルーティングを行う
func SubdomainRouter(w http.ResponseWriter, r *http.Request) {
	subdomainCtx := r.Context().Value("subdomain").(*SubdomainContext)

	if subdomainCtx == nil {
		http.Error(w, "Invalid subdomain", http.StatusBadRequest)
		return
	}

	switch subdomainCtx.Subdomain {
	case "mysite7":
		// mysite7専用の処理
		handleMysite7(w, r)
	case "mysite2":
		// mysite2専用の処理
		handleMysite2(w, r)
	default:
		// デフォルトサイト処理
		handleDefaultSite(w, r)
	}
}

// handleMysite7 はmysite7専用の処理
func handleMysite7(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/polls") {
		// polls機能の処理
		servePollsPage(w, r)
	} else {
		// mysite7のメインページ
		serveMysite7Home(w, r)
	}
}

// handleMysite2 はmysite2専用の処理
func handleMysite2(w http.ResponseWriter, r *http.Request) {
	// mysite2専用のロジック
	serveMysite2(w, r)
}

// handleDefaultSite はデフォルトサイトの処理
func handleDefaultSite(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Site not found", http.StatusNotFound)
}

// 各サイト専用のハンドラー関数（実装例）
func servePollsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Polls - MySite7</title>
</head>
<body>
    <h1>投票システム</h1>
    <p>mysite7.fumi042-server.top の投票ページです</p>
</body>
</html>
	`))
}

func serveMysite7Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>MyS	ite7</title>
</head>
<body>
    <h1>MyS	ite7 ホーム</h1>
    <p><a href="/polls">投票ページへ</a></p>
</body>
</html>
	`))
}

func serveMysite2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head>
    <title>MyS	ite2</title>
</head>
<body>
    <h1>MyS	ite2</h1>
    <p>mysite2.fumi042-server.top のコンテンツです</p>
</body>
</html>
	`))
}
