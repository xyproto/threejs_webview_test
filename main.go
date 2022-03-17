package main

import (
	"net/http"

	"github.com/urfave/negroni"
	"github.com/xyproto/onthefly"
	"github.com/webview/webview"
)

// Create a Three.JS page
func ThreeJSPage() *onthefly.Page {
	p, t := onthefly.NewThreeJSWithGiven("Embedded Three.js", threejsSrc)

	// Add a camera at (0, 0, 5)
	t.AddCamera()
	t.CameraPos("z", 5)

	// We also need a renderer
	t.AddRenderer()

	// Create a test cube
	cube1 := t.AddTestCube()

	// Create another test cube, and rotate it a bit
	cube2 := t.AddTestCube()
	cube2.JS += cube2.ID + ".rotation.x += 0.9;"

	// Render function (happens every frame)
	r := onthefly.NewRenderFunction()

	// Rotate the first cube
	r.AddJS(cube1.ID + ".rotation.x += 0.02;")
	r.AddJS(cube1.ID + ".rotation.y += 0.02;")

	// Rotate the second cube at a different speed
	r.AddJS(cube2.ID + ".rotation.x += 0.04;")
	r.AddJS(cube2.ID + ".rotation.y += 0.07;")

	// Add the render function to the script tag
	t.AddRenderFunction(r, true)

	return p
}

func main() {

	// Create a Negroni instance and a ServeMux instance
	n := negroni.Classic()
	mux := http.NewServeMux()

	// Create the page by calling the function above
	page := ThreeJSPage()

	// Publish the generated page (HTML and CSS)
	page.Publish(mux, "/", "/style.css", false)

	// Handler goes last
	n.UseHandler(mux)

	// Listen for requests at port 1814
	go func() {
		n.Run(":1814")
	}()

	// Open three.js 3D graphics in a 1024x768 resizable window
	debug := true
	wv := webview.New(debug)
	defer wv.Destroy()
	wv.SetTitle("Three JS")
	wv.SetSize(1024, 768, webview.HintNone)
	wv.Navigate("http://localhost:1814/")
	wv.Run()
}
