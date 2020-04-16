package assets

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetse8974eaed9595d681b6ed62f7f659dc0fb930889 = "<form action=\"/devices\" method=\"post\" class=\"row\">\n  <div class=\"col-md-6\">\n    <div class=\"form-group\">\n      <label>Device Name</label>\n      <input type=\"text\" name=\"name\" class=\"form-control\" value=\"{{ .newDeviceName }}\" placeholder=\"ex: Work Laptop\" />\n    </div>\n    <div class=\"form-group\">\n      <label>Platform</label>\n      <select name=\"os\" class=\"form-control\">\n        <option selected disabled value=\"\">Please select</option>\n        <option value=\"mac\">Mac</option>\n        <option value=\"linux\">Linux</option>\n        <option value=\"windows\">Windows</option>\n        <option value=\"ios\">iOS</option>\n        <option value=\"android\">Android</option>\n      </select>\n    </div>\n    <button type=\"submit\" class=\"btn btn-primary mb-2\">Add Device</button>\n  </div>\n</form>"
var _Assetsf441836287dd10b66acc92cc115787cea88b1fac = "{{ template \"layout_header.html\" }}\n\n<div class=\"text-center\">\n  <h1>Authentication</h1>\n  <p>Please sign in with your Google account to continue</p>\n  <div>\n    <a href=\"{{ .startPath }}\" class=\"btn btn-primary\">\n      <i class=\"fab fa-google\"></i> Sign in with Google\n    </a>\n  </div>\n</div>\n\n{{ template \"layout_footer.html\" }}"
var _Assetsdc752b78544834f324e9bd47f84b02df0f410ac8 = "{{ template \"layout_header.html\" }}\n\n<div class=\"block\">\n  <div class=\"card\">\n    <div class=\"card-header\">\n      Server Details\n    </div>\n    <div class=\"card-body\">\n      <dl class=\"row\">\n        <dt class=\"col-sm-3\">Name</dt>\n        <dd class=\"col-sm-9\">{{ .server.Name }}</dd>\n      </dl>\n      <dl class=\"row\">\n        <dt class=\"col-sm-3\">Interface</dt>\n        <dd class=\"col-sm-9\">\n          {{ .server.Interface }} - {{ .server.IPV4Net }} <span class=\"text-muted\">({{ .server.IPV4Count }} addresses)</span>\n        </dd>\n      </dl>\n      <dl class=\"row\">\n        <dt class=\"col-sm-3\">Endpoint</dt>\n        <dd class=\"col-sm-9\">{{ .server.PublicAddr }}</dd>\n      </dl>\n      {{ if ne .server.DNS \"\" }}\n        <dl class=\"row\">\n          <dt class=\"col-sm-3\">DNS</dt>\n          <dd class=\"col-sm-9\">{{ .server.DNS }}</dd>\n        </dl>\n      {{ end }}\n      <dl class=\"row\">\n        <dt class=\"col-sm-3\"></dt>\n        <dd class=\"col-sm9\">\n          <a href=\"/admin/server\" class=\"btn btn-outline-primary\">Settings</a>\n          <a href=\"/admin/server/restart\" class=\"btn btn-outline-secondary\">Reload Configuration</a>\n        </dd>\n      </dl>\n    </div>\n  </div>\n</div>\n\n<div class=\"block\">\n  <h1>Users</h1>\n  <table class=\"table table-bordered\">\n    <thead>\n      <tr>\n        <th>ID</th>\n        <th>Name</th>\n        <th>Email</th>\n        <th>Role</th>\n        <th>Created</th>\n      </tr>\n    </thead>\n    <tbody>\n      {{ range .users }}\n      <tr>\n        <td>{{ .ID }}</td>\n        <td>{{ .Name }}</td>\n        <td>{{ .Email }}</td>\n        <td>{{ .Role }}</td>\n        <td>{{ time .CreatedAt }}</td>\n      </tr>\n      {{ end }}\n    </tbody>\n  </table>\n</div>\n\n<div class=\"block\">\n  <h1>Devices</h1>\n  <table class=\"table table-bordered\">\n    <thead>\n      <tr>\n        <th>ID</th>\n        <th>User</th>\n        <th>Name</th>\n        <th>IP</th>\n        <th>Created</th>\n        <th>Last Seen</th>\n        <th>Endpoint</th>\n      </tr>\n    </thead>\n    <tbody>\n      {{ range .devices }}\n        <tr>\n          <td>{{ .ID }}</td>\n          <td>{{ .UserID }}</td>\n          <td>{{ .Name }}</td>\n          <td>{{ .IPV4 }}</td>\n          <td>{{ time .CreatedAt }}</td>\n          {{ with .GetPeerInfo }}\n            <td>{{ time .LastHandshakeTime }}</td>\n            <td>{{ with .Endpoint }}{{ . }}{{ end }}</td>\n          {{ else }}\n            <td colspan=\"2\">\n              <span class=\"text-muted\">N/A</span>\n            </td>\n          {{ end }}\n        </tr>\n      {{ end }}\n    </tbody>\n  </table>\n</div>\n\n<div class=\"block\">\n  <h1>Wireguard Config</h1>\n  <div class=\"card card-body\">\n    <pre>{{ .config }}</pre>\n  </div>\n</div>\n\n{{ template \"layout_footer.html\" }}"
var _Assetsd3607537e7e9216275a956ff154746d7eb53630d = "</div></div></div>\n</body>\n</html>\n"
var _Assetsd9961776eee0ed6e9bd47203fcb44195f76ed068 = "{{ template \"layout_header.html\" }}\n\n<div class=\"text-center\">\n  <h2>How to Connect</h2>\n  <p>Follow the instructions below to setup your WireGuard connection.</p>\n</div>\n\n<div class=\"card device\">\n  <div class=\"card-header\">\n    {{ .device.Name }}\n  </div>\n  <div class=\"card-body\">\n    <ul>\n      <li>\n        <a href=\"https://www.wireguard.com/install/\" target=\"_blank\">Download WireGuard</a> client for your operating system.\n      </li>\n      {{ if .device.IsMobile }}\n        <li>\n          <div>Open WireGuard application and add a new tunnel by scanning the barcode below.</div>\n          <div>Alternatively, you can download the config and add it manually.</div>\n          <div class=\"text-center\">\n            <img src=\"/devices/{{ .device.ID }}/config?qr=1\" class=\"qr\" />\n          </div>\n        </li>\n      {{ else }}\n        <li>\n          Download the tunnel configuration file by clicking on the button below.\n        </li>\n        <li>\n          Open WireGuard client application and add a new tunnel from the downloaded file.\n        </li>  \n      {{ end }}\n    </ul>\n\n    <div class=\"config-download\">\n      <a href=\"/devices/{{ .device.ID }}/config\" class=\"btn btn-outline-primary\">Download WireGuard Config</a> \n      <a href=\"/\" class=\"btn btn-outline-secondary\">Go Back</a><br />\n    </div>\n  </div>\n</div>\n\n{{ template \"layout_footer.html\" }}"
var _Assets9d8eb98e0454ca4aa1fae7dd60913a2c389ca33c = "{{ template \"layout_header.html\" }}\n\n<div class=\"alert alert-danger\">\n  <h4 class=\"alert-heading\">Error has occurred!</h4>\n  <p>{{ .error }}</p>\n  <div><a href=\"javascript:history.back()\" class=\"btn btn-danger\">Go Back</a></div>\n</div>\n\n{{ template \"layout_footer.html\" }}"
var _Assets8ea2097b9a9db97de929c0dd70f549f60e8d1b4c = "{{ template \"layout_header.html\" }}\n\n<div class=\"card block account\">\n  <div class=\"card-body\">\n    <div class=\"row\">\n      <div class=\"col-md-9\">\n        You are signed in as <u>{{ .user.Email }}</u>\n      </div>\n      <div class=\"col-md-3 text-right\">\n        {{ if eq .user.Role \"admin\" }}\n          <a href=\"/admin\" class=\"btn btn-outline-primary btn-sm\">Admin</a>\n        {{ end }}\n        <a href=\"/auth/signout\" class=\"btn btn-outline-secondary btn-sm\">Sign Out</a>\n      </div>\n    </div>\n  </div>\n</div>\n\n<h3>Your Devices</h3>\n<table class=\"table devices\">\n  <tbody>\n    {{ range .devices }}\n    <tr>\n      <td>\n        <i class=\"{{ .IconClass }}\"></i>\n        {{ .Name }}\n      </td>\n      <td>{{ .IPV4 }}</td>\n      <td class=\"text-right\">\n        <a href=\"/devices/{{ .ID }}\" class=\"btn btn-outline-secondary btn-sm\">Show Config</a>\n        <a href=\"/devices/{{ .ID }}/delete\" class=\"btn btn-outline-danger btn-sm\">Remove</a>\n      </td>\n    </tr>\n    {{ end }}\n  </tbody>\n</table>\n\n<h3>Add a new device</h3>\n{{ template \"device_form.html\" }}\n\n{{ template \"layout_footer.html\" }}"
var _Assetsec699bf9721f88369446cff17880caf54809e2cc = "{{ template \"layout_header.html\" }}\n\n<h1>Configure WireGuard server</h1>\n\n{{ if .error }}\n  <div class=\"alert alert-danger\">{{ .error }}</div>\n{{ end }}\n\n<form action=\"/admin/server\" method=\"post\">\n  <div class=\"form-group\">\n    <label>Name</label>\n    <input class=\"form-control\" type=\"text\" name=\"name\" value=\"{{ .server.Name }}\" />\n  </div>\n  <div class=\"row\">\n    <div class=\"form-group col-md-8\">\n      <label>Endpoint</label>\n      <input class=\"form-control\" type=\"text\" name=\"endpoint\" value=\"{{ .server.Endpoint }}\" />\n    </div>\n    <div class=\"form-group col-md-4\">\n      <label>Listen Port</label>\n      <input class=\"form-control\" type=\"text\" name=\"listen_port\" value=\"{{ .server.ListenPort }}\" />\n    </div>\n  </div>\n  <div class=\"row\">\n    <div class=\"form-group col-md-4\">\n      <label>Interface</label>\n      <input class=\"form-control\" type=\"text\" name=\"interface\" value=\"{{ .server.Interface }}\" />\n    </div>\n    <div class=\"form-group col-md-4\">\n      <label>Network</label>\n      <input class=\"form-control\" type=\"text\" name=\"ipv4_net\" value=\"{{ .server.IPV4Net }}\" />\n    </div>\n    <div class=\"form-group col-md-4\">\n      <label>Connect Mode</label>\n      <select name=\"mode\" class=\"form-control\">\n        <option value=\"split\" selected>Split Tunnel</option>\n        <option value=\"full\">Full Tunnel</option>\n      </select>\n    </div>\n  </div>\n  <div class=\"form-group\">\n    <label>DNS Servers</label>\n    <input class=\"form-control\" type=\"text\" name=\"dns\" value=\"{{ .server.DNS }}\" placeholder=\"Optional\" />\n  </div>\n  <div class=\"form-group\">\n    <label>Post-Up Commands</label>\n    <input class=\"form-control\" type=\"text\" name=\"interface\" value=\"{{ .server.PostUp }}\" />\n  </div>\n  <div class=\"form-group\">\n    <label>Post-Down Commands</label>\n    <input class=\"form-control\" type=\"text\" name=\"interface\" value=\"{{ .server.PostDown }}\" />\n  </div>\n  <div class=\"form-actions\">\n    <button type=\"submit\" class=\"btn btn-primary\">Save Changes</button>\n  </div>\n</form>\n\n{{ template \"layout_footer.html\" }}"
var _Assets51a8d24399ae4ab763f47c5bd865fad3b2b7200c = "body {\n  padding: 50px 0px;\n}\n\n.block {\n  margin-bottom: 20px;\n}\n\n.card.device .card-body {\n  padding: 12px;\n}\n\n.qr {\n  margin: 12px 0px;\n  border: 2px solid #eee;\n  border-radius: 4px;;\n  background: #ff00ff;\n}\n\n.config-download {\n  margin: 12px 0px;\n  text-align: center;\n}"
var _Assets53588d969ca4aa48419a4d621aae4391b3db7056 = "<html>\n  <head>\n    <link rel=\"stylesheet\" href=\"https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css\" integrity=\"sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T\" crossorigin=\"anonymous\">\n    <link rel=\"stylesheet\" href=\"/static/main.css\">\n    <script src=\"https://kit.fontawesome.com/e654a3ca3a.js\" crossorigin=\"anonymous\"></script>\n  </head>\n  <body>\n    <div class=\"container\">\n      <div class=\"row\">\n        <div class=\"col-md-12\">\n      "

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"static"}, "/static": []string{"index.html", "admin_server.html", "device_form.html", "login.html", "main.css", "layout_header.html", "admin_index.html", "layout_footer.html", "device.html", "error.html"}}, map[string]*assets.File{
	"/static": &assets.File{
		Path:     "/static",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1587007287, 1587007287885876026),
		Data:     nil,
	}, "/static/device_form.html": &assets.File{
		Path:     "/static/device_form.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586668809, 1586668809304370915),
		Data:     []byte(_Assetse8974eaed9595d681b6ed62f7f659dc0fb930889),
	}, "/static/login.html": &assets.File{
		Path:     "/static/login.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586641384, 1586641384194646128),
		Data:     []byte(_Assetsf441836287dd10b66acc92cc115787cea88b1fac),
	}, "/static/admin_index.html": &assets.File{
		Path:     "/static/admin_index.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1587011362, 1587011362306505246),
		Data:     []byte(_Assetsdc752b78544834f324e9bd47f84b02df0f410ac8),
	}, "/static/layout_footer.html": &assets.File{
		Path:     "/static/layout_footer.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586641384, 1586641384194424554),
		Data:     []byte(_Assetsd3607537e7e9216275a956ff154746d7eb53630d),
	}, "/static/device.html": &assets.File{
		Path:     "/static/device.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586641384, 1586641384193974450),
		Data:     []byte(_Assetsd9961776eee0ed6e9bd47203fcb44195f76ed068),
	}, "/static/error.html": &assets.File{
		Path:     "/static/error.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586641384, 1586641384194191236),
		Data:     []byte(_Assets9d8eb98e0454ca4aa1fae7dd60913a2c389ca33c),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1587011416, 1587011416096539639),
		Data:     nil,
	}, "/static/index.html": &assets.File{
		Path:     "/static/index.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586668669, 1586668669348694814),
		Data:     []byte(_Assets8ea2097b9a9db97de929c0dd70f549f60e8d1b4c),
	}, "/static/admin_server.html": &assets.File{
		Path:     "/static/admin_server.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1587009737, 1587009737160082192),
		Data:     []byte(_Assetsec699bf9721f88369446cff17880caf54809e2cc),
	}, "/static/main.css": &assets.File{
		Path:     "/static/main.css",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586641384, 1586641384194812246),
		Data:     []byte(_Assets51a8d24399ae4ab763f47c5bd865fad3b2b7200c),
	}, "/static/layout_header.html": &assets.File{
		Path:     "/static/layout_header.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1587010145, 1587010145202435326),
		Data:     []byte(_Assets53588d969ca4aa48419a4d621aae4391b3db7056),
	}}, "")
