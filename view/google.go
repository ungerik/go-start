package view

import (
	//	"os"
	"fmt"
	"net/url"
	//	"github.com/ungerik/go-start/utils"
)

func GoogleAnalytics(trackingID string) HTML {
	return Printf("<script>window._gaq = [['_setAccount','%s'],['_trackPageview'],['_trackPageLoadTime']];Modernizr.load({load: ('https:' == location.protocol ? '//ssl' : '//www') + '.google-analytics.com/ga.js'});</script>", trackingID)
}

// todo: replace http with https if necessary
func GoogleMaps(apiKey string, sensor bool, callback string) HTML {
	return Printf("<script src='http://maps.googleapis.com/maps/api/js?key=%s&sensor=%t&callback=%s'></script>", apiKey, sensor, callback)
}

func GoogleMapsIframe(width, height int, location string) *Iframe {
	location = url.QueryEscape(location)
	URL := fmt.Sprintf("http://maps.google.com/maps?q=%s&output=embed", location)
	//URL = url.QueryEscape(URL)
	return &Iframe{
		Width:  width,
		Height: height,
		URL:    URL,
	}
}

//type GoogleMapType string
//
//const (
//	GoogleMapTypeHybrid    GoogleMapType = "google.maps.MapTypeId.HYBRID"    //This map type displays a transparent layer of major streets on satellite images.
//	GoogleMapTypeRoadmap   GoogleMapType = "google.maps.MapTypeId.ROADMAP"   //This map type displays a normal street map.
//	GoogleMapTypeSatellite GoogleMapType = "google.maps.MapTypeId.SATELLITE" //This map type displays satellite images.
//	GoogleMapTypeTerrain   GoogleMapType = "google.maps.MapTypeId.TERRAIN"   //This map type displays maps with physical features such as terrain and vegetation.
//)
//
/////////////////////////////////////////////////////////////////////////////////
//// GoogleMap
//
//type GoogleMap struct {
//	ViewBaseWithId
//	Class     string
//	Width     int
//	Height    int
//	Type      GoogleMapType
//	CenterLat float64
//	CenterLng float64
//	Zoom      float64
//}
//
//func (self *GoogleMap) Render(ctx *Context) (err error) {
//	return nil
//}
