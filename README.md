> [!TIP]
> View this project live at [go.baylor.dev](https://go.baylor.dev/)!

# 💫 Go Vanity ✨
Go Vanity is a simple web server that allows you to host Go modules with your own vanity url. For example, instead of ...
```
go get github.com/baylorhenshaw/go-vanity
```
You can use ...
```
go get go.baylor.dev/go-vanity
```
# 💾 Installation 💿

> [!WARNING]
> Installation instructions coming soon!

# 🔧 Configuration 🔨

### Modules File

> [!IMPORTANT]
> You must define the path for this file in your environment variables!

In order to configure Go Vanity, you need a `.json` file containing your module configuration. Here is an example ...
```json
{
    "modules": [
        {
            "name": "go-vanity",
            "github": "https://github.com/baylorhenshaw/go-vanity"
        }
    ]
}
```
This will add the vanity url `https://yourdomain.com/go-vanity/` and it will redirect to the proper Github repository. You can add as many modules as you need in this config file.

### Environment Variables

You can get started with Go Vanity with these environment variables ...

```sh
HOST=yourdomain.com
PORT=3000
MODULES_FILE=modules.json
```

### Umami Analytics

You can enable Umami Analytics by adding the following environment variables ...

```sh
UMAMI_URL=<insert-umami-url>
UMAMI_ID=<insert-website-id>
```
