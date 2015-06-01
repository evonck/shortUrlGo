package controllers

import "github.com/revel/revel"
import "github.com/evonck/shortlink/mylib"
import "hash/fnv"
import "net/http"
import "net/url"
import "strconv"
import "time"



type App struct {
	*revel.Controller
	mylib.RedisController
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) ShortenUrl() revel.Result {
	var longUrl string = c.Params.Get("longUrl")
	CheckIfValidUrl(c,longUrl)
	var preference string = c.Params.Get("preferences")
	var shortUrlCode string
	if(preference == ""){
    	shortUrlCode = strconv.FormatUint(uint64( hash(longUrl)),10)
	}else{
    	shortUrlCode = preference
	}
	shortUrl := CreateUrl(shortUrlCode)
	LinkAlreadyExist, err := c.Redis.Get(shortUrl).Result()
	if (err != nil) {
        StroeInRedis(c,shortUrl,longUrl)
		return c.RenderJson(shortUrl)
    }else if (LinkAlreadyExist == longUrl) {
		return c.RenderJson(shortUrl)
    }else if (preference != ""){
		var CountPreference int64 = 0
    	for (LinkAlreadyExist != ""){
			shortUrlCode = preference + strconv.FormatInt(CountPreference,10);   
			shortUrl = CreateUrl(shortUrlCode)
			LinkAlreadyExist,err = c.Redis.Get(shortUrl).Result()
			CountPreference++;
		}
		return c.RenderJson(shortUrl)
    }else{
    	c.Response.Status = 409
		return c.RenderJson("The ShortUrlCode already correspond to another Long URL.")
    }
	return c.RenderJson(shortUrl)

}

func (c App) GetLongUrl(shortUrl string) revel.Result {
	var longUrl string = c.Params.Get("shortUrl")
	longUrl, err := c.Redis.Get(shortUrl).Result()
    if (err != nil) {
        c.Response.Status = 200
		return c.RenderJson("does not exists")
    }
	return c.RenderJson(longUrl)
}

func hash(s string) uint32 {
        h := fnv.New32a()
        h.Write([]byte(s))
        return h.Sum32()
}

func CreateUrl(shortUrlCode string) string {
	domainAddr, found := revel.Config.String("domaineName")
	if !found {
		revel.ERROR.Fatal("No domaineName found.")
	}
    shortUrl := &url.URL{
		Host:   domainAddr,
		Scheme: "http",
		Path: shortUrlCode,
	}	
	return shortUrl.String()
}

func CheckIfValidUrl(c App,url string){
	resp, err := http.Get("http://"+url)
	if err != nil {
		revel.ERROR.Fatal("False URL")
	}
    defer resp.Body.Close()
}

func StroeInRedis(c App,shortUrl string,longUrl string){
	err := c.Redis.Set(shortUrl,longUrl,7889231*time.Second).Err()
    if err != nil {
		revel.ERROR.Fatal("Error while storing in db")
    }
}

