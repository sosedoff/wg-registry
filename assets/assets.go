package assets

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets8ea2097b9a9db97de929c0dd70f549f60e8d1b4c = "{{ template \"layout_header.html\" }}\n\n<div class=\"card block account\">\n  <div class=\"card-body\">\n    <div class=\"row\">\n      <div class=\"col-md-9\">\n        You are signed in as <u>{{ .user.Email }}</u>\n      </div>\n      <div class=\"col-md-3 text-right\">\n        {{ if eq .user.Role \"admin\" }}\n          <a href=\"/admin\" class=\"btn btn-outline-primary btn-sm\">Admin</a>\n        {{ end }}\n        <a href=\"/auth/signout\" class=\"btn btn-outline-secondary btn-sm\">Sign Out</a>\n      </div>\n    </div>\n  </div>\n</div>\n\n<h3>Your Devices</h3>\n<table class=\"table devices\">\n  <tbody>\n    {{ range .devices }}\n    <tr>\n      <td>\n        <i class=\"{{ .IconClass }}\"></i>\n        {{ .Name }}\n      </td>\n      <td>{{ .IPV4 }}</td>\n      <td class=\"text-right\">\n        <a href=\"/devices/{{ .ID }}\" class=\"btn btn-outline-secondary btn-sm\">Show Config</a>\n        <a href=\"/devices/{{ .ID }}/delete\" class=\"btn btn-outline-danger btn-sm\">Remove</a>\n      </td>\n    </tr>\n    {{ end }}\n  </tbody>\n</table>\n\n<h3>Add a new device</h3>\n{{ template \"device_form.html\" }}\n\n{{ template \"layout_footer.html\" }}"
var _Assets51a8d24399ae4ab763f47c5bd865fad3b2b7200c = "body {\n  padding: 50px 0px;\n}\n\n.block {\n  margin-bottom: 20px;\n}\n\n.card.device .card-body {\n  padding: 12px;\n}\n\n.qr {\n  margin: 12px 0px;\n  border: 2px solid #eee;\n  border-radius: 4px;;\n  background: #ff00ff;\n}\n\n.config-download {\n  margin: 12px 0px;\n  text-align: center;\n}"
var _Assets1fa295fb8b4dabc19ec65086f2b0da02d4a18545 = "{{ template \"layout_header.html\" }}\n\n<div class=\"block\">\n  <h1>Peers</h1>\n  <table class=\"table table-bordered\">\n    <thead>\n      <tr>\n        <th>Key</th>\n        <th>Endpoint</th>\n        <th>IPs</th>\n        <th>Handshake</th>\n        <th>Sent</th>\n        <th>Received</th>\n      </tr>\n    </thead>\n    <tbody>\n      {{ range .peers }}\n      <tr>\n        <td>{{ .PublicKey }}</td>\n        <td>{{ .Endpoint }}</td>\n        <td>{{ .AllowedIPs }}</td>\n        <td>{{ .LastHandshakeTime }}</td>\n        <td>{{ .TransmitBytes }}</td>\n        <td>{{ .ReceiveBytes }}</td>\n      </tr>\n      {{ end }}\n    </tbody>\n  </table>\n</div>\n\n{{ template \"layout_footer.html\" }}"
var _Assetsdc752b78544834f324e9bd47f84b02df0f410ac8 = "{{ template \"layout_header.html\" }}\n\n<div class=\"block\">\n  <h1>Server</h1>\n  <div class=\"card card-body\">\n    <dl class=\"row\">\n      <dt class=\"col-sm-3\">Name</dt>\n      <dd class=\"col-sm-9\">{{ .server.Name }}</dd>\n    </dl>\n    <dl class=\"row\">\n      <dt class=\"col-sm-3\">Interface</dt>\n      <dd class=\"col-sm-9\">\n        {{ .server.Interface }} - {{ .server.IPV4Net }}\n        <a href=\"/admin/peers\" class=\"btn btn-outline-secondary btn-sm\">Show Peers</a>\n      </dd>\n    </dl>\n    <dl class=\"row\">\n      <dt class=\"col-sm-3\">Endpoint</dt>\n      <dd class=\"col-sm-9\">{{ .server.PublicAddr }}</dd>\n    </dl>\n    <dl class=\"row\">\n      <dt class=\"col-sm-3\">DNS</dt>\n      <dd class=\"col-sm-9\">{{ .server.DNS }}</dd>\n    </dl>\n  </div>\n</div>\n\n<div class=\"block\">\n  <h1>Users</h1>\n  <table class=\"table table-bordered\">\n    <thead>\n      <tr>\n        <th>ID</th>\n        <th>Name</th>\n        <th>Email</th>\n        <th>Role</th>\n        <th>Created</th>\n      </tr>\n    </thead>\n    <tbody>\n      {{ range .users }}\n      <tr>\n        <td>{{ .ID }}</td>\n        <td>{{ .Name }}</td>\n        <td>{{ .Email }}</td>\n        <td>{{ .Role }}</td>\n        <td>{{ time .CreatedAt }}</td>\n      </tr>\n      {{ end }}\n    </tbody>\n  </table>\n</div>\n\n<div class=\"block\">\n  <h1>Devices</h1>\n  <table class=\"table table-bordered\">\n    <thead>\n      <tr>\n        <th>ID</th>\n        <th>User</th>\n        <th>Name</th>\n        <th>IP</th>\n        <th>Active</th>\n        <th>Created</th>\n      </tr>\n    </thead>\n    <tbody>\n      {{ range .devices }}\n        <tr>\n          <td>{{ .ID }}</td>\n          <td>{{ .UserID }}</td>\n          <td>{{ .Name }}</td>\n          <td>{{ .IPV4 }}</td>\n          <td>{{ .Enabled }}</td>\n          <td>{{ time .CreatedAt }}</td>\n        </tr>\n      {{ end }}\n    </tbody>\n  </table>\n</div>\n\n<div class=\"block\">\n  <h1>Wireguard Config</h1>\n  <div class=\"card card-body\">\n    <pre>{{ .config }}</pre>\n  </div>\n</div>\n\n{{ template \"layout_footer.html\" }}"
var _Assetse8974eaed9595d681b6ed62f7f659dc0fb930889 = "<form action=\"/devices\" method=\"post\">\n  <div class=\"form-group mb-2\">\n    <label class=\"my-1 mr-2\">Name</label>\n    <input type=\"text\" name=\"name\" class=\"form-control\" value=\"{{ .newDeviceName }}\" placeholder=\"ex: Work Laptop\" />\n  </div>\n  <div class=\"form-group mb-2 text-small\">\n    <label>Platform</label>\n    <select name=\"os\" class=\"form-control\">\n      <option selected disabled value=\"\">Please select</option>\n      <option value=\"mac\">Mac</option>\n      <option value=\"linux\">Linux</option>\n      <option value=\"windows\">Windows</option>\n      <option value=\"ios\">iOS</option>\n      <option value=\"android\">Android</option>\n    </select>\n  </div>\n  <button type=\"submit\" class=\"btn btn-primary mb-2\">Add Device</button>\n</form>"
var _Assetsf441836287dd10b66acc92cc115787cea88b1fac = "{{ template \"layout_header.html\" }}\n\n<div class=\"text-center\">\n  <h1>Authentication</h1>\n  <p>Please sign in with your Google account to continue</p>\n  <div>\n    <a href=\"{{ .startPath }}\" class=\"btn btn-primary\">\n      <i class=\"fab fa-google\"></i> Sign in with Google\n    </a>\n  </div>\n</div>\n\n{{ template \"layout_footer.html\" }}"
var _Assets53588d969ca4aa48419a4d621aae4391b3db7056 = "<html>\n  <head>\n    <link rel=\"stylesheet\" href=\"https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css\" integrity=\"sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T\" crossorigin=\"anonymous\">\n    <link rel=\"stylesheet\" href=\"/static/main.css\">\n    <script src=\"https://cdnjs.cloudflare.com/ajax/libs/jquery/3.4.0/jquery.min.js\"></script>\n    <script src=\"https://kit.fontawesome.com/e654a3ca3a.js\" crossorigin=\"anonymous\"></script>\n  </head>\n  <body>\n    <div class=\"container\">\n      <div class=\"row\">\n        <div class=\"col-md-8 offset-md-2\">\n      "
var _Assetsd3607537e7e9216275a956ff154746d7eb53630d = "</div></div></div>\n</body>\n</html>\n"
var _Assetsd9961776eee0ed6e9bd47203fcb44195f76ed068 = "{{ template \"layout_header.html\" }}\n\n<div class=\"text-center\">\n  <h2>How to Connect</h2>\n  <p>Follow the instructions below to setup your WireGuard connection.</p>\n</div>\n\n<div class=\"card device\">\n  <div class=\"card-header\">\n    {{ .device.Name }}\n  </div>\n  <div class=\"card-body\">\n    <ul>\n      <li>\n        <a href=\"https://www.wireguard.com/install/\" target=\"_blank\">Download WireGuard</a> client for your operating system.\n      </li>\n      {{ if .device.IsMobile }}\n        <li>\n          <div>Open WireGuard application and add a new tunnel by scanning the barcode below.</div>\n          <div>Alternatively, you can download the config and add it manually.</div>\n          <div class=\"text-center\">\n            <img src=\"/devices/{{ .device.ID }}/config?qr=1\" class=\"qr\" />\n          </div>\n        </li>\n      {{ else }}\n        <li>\n          Download the tunnel configuration file by clicking on the button below.\n        </li>\n        <li>\n          Open WireGuard client application and add a new tunnel from the downloaded file.\n        </li>  \n      {{ end }}\n    </ul>\n\n    <div class=\"config-download\">\n      <a href=\"/devices/{{ .device.ID }}/config\" class=\"btn btn-outline-primary\">Download WireGuard Config</a> \n      <a href=\"/\" class=\"btn btn-outline-secondary\">Go Back</a><br />\n    </div>\n  </div>\n</div>\n\n{{ template \"layout_footer.html\" }}"
var _Assets9d8eb98e0454ca4aa1fae7dd60913a2c389ca33c = "{{ template \"layout_header.html\" }}\n\n<div class=\"alert alert-danger\">\n  <h4 class=\"alert-heading\">Error has occurred!</h4>\n  <p>{{ .error }}</p>\n  <div><a href=\"javascript:history.back()\" class=\"btn btn-danger\">Go Back</a></div>\n</div>\n\n{{ template \"layout_footer.html\" }}"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"static"}, "/static": []string{"index.html", "device_form.html", "login.html", "main.css", "layout_header.html", "admin_peers.html", "admin_index.html", "layout_footer.html", "device.html", "error.html"}}, map[string]*assets.File{
	"/static/layout_header.html": &assets.File{
		Path:     "/static/layout_header.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586637620, 1586637620132040908),
		Data:     []byte(_Assets53588d969ca4aa48419a4d621aae4391b3db7056),
	}, "/static/layout_footer.html": &assets.File{
		Path:     "/static/layout_footer.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586227077, 1586227077387926034),
		Data:     []byte(_Assetsd3607537e7e9216275a956ff154746d7eb53630d),
	}, "/static/device.html": &assets.File{
		Path:     "/static/device.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586575640, 1586575640612917939),
		Data:     []byte(_Assetsd9961776eee0ed6e9bd47203fcb44195f76ed068),
	}, "/static/error.html": &assets.File{
		Path:     "/static/error.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586639150, 1586639150440957235),
		Data:     []byte(_Assets9d8eb98e0454ca4aa1fae7dd60913a2c389ca33c),
	}, "/static/device_form.html": &assets.File{
		Path:     "/static/device_form.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586570041, 1586570041055000000),
		Data:     []byte(_Assetse8974eaed9595d681b6ed62f7f659dc0fb930889),
	}, "/static/login.html": &assets.File{
		Path:     "/static/login.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586227093, 1586227093177336368),
		Data:     []byte(_Assetsf441836287dd10b66acc92cc115787cea88b1fac),
	}, "/static/index.html": &assets.File{
		Path:     "/static/index.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586637681, 1586637681197567740),
		Data:     []byte(_Assets8ea2097b9a9db97de929c0dd70f549f60e8d1b4c),
	}, "/static/main.css": &assets.File{
		Path:     "/static/main.css",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586565481, 1586565481730916055),
		Data:     []byte(_Assets51a8d24399ae4ab763f47c5bd865fad3b2b7200c),
	}, "/static/admin_peers.html": &assets.File{
		Path:     "/static/admin_peers.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586638840, 1586638840511328710),
		Data:     []byte(_Assets1fa295fb8b4dabc19ec65086f2b0da02d4a18545),
	}, "/static/admin_index.html": &assets.File{
		Path:     "/static/admin_index.html",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1586638912, 1586638912608621866),
		Data:     []byte(_Assetsdc752b78544834f324e9bd47f84b02df0f410ac8),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1586627973, 1586627973824116266),
		Data:     nil,
	}, "/static": &assets.File{
		Path:     "/static",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1586638753, 1586638753459490596),
		Data:     nil,
	}}, "")
