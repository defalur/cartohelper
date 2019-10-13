package mapWidget

import (
    "fmt"
    "fyne.io/fyne"
    "fyne.io/fyne/widget"
    "fyne.io/fyne/canvas"
    "image"
    "image/color"
    "cartohelper/mapViewer"
    "os"
    "image/png"
)

type MapWidget struct {
    size fyne.Size
    position fyne.Position
    hidden bool
    clickx int
    clicky int
    Scale float32
    MapViewer mapviewer.MapViewer
    lastImage *image.RGBA
}

func (m *MapWidget) SaveImg(name string) {
    out, _ := os.Create(name + ".png")
    
    png.Encode(out, m.lastImage)
}

func (m *MapWidget) Size() fyne.Size {
    return m.size
}

func (m *MapWidget) Resize(size fyne.Size) {
    m.size = size
    widget.Renderer(m).Layout(size)
}

func (m *MapWidget) Position() fyne.Position {
    return m.position
}

func (m *MapWidget) Move(pos fyne.Position) {
    m.position = pos
    widget.Renderer(m).Layout(m.size)
}

func (m *MapWidget) MinSize() fyne.Size {
    return widget.Renderer(m).MinSize()
}

func (m *MapWidget) Visible() bool {
    return !m.hidden
}

func (m *MapWidget) Show() {
    m.hidden = false
}

func (m *MapWidget) Hide() {
    m.hidden = true
}

type MapWidgetRenderer struct {
    render *canvas.Raster
    objects []fyne.CanvasObject
    
    state *MapWidget
}

func (r *MapWidgetRenderer) MinSize() fyne.Size {
    return fyne.NewSize(r.state.MapViewer.MapState().GetWidth(),
                        r.state.MapViewer.MapState().GetHeight())
}

func (r *MapWidgetRenderer) Layout(size fyne.Size) {
    r.render.Resize(size)
}

func (r *MapWidgetRenderer) ApplyTheme() {
    
}

func (r *MapWidgetRenderer) BackgroundColor() color.Color {
    return color.RGBA{128, 128, 128, 255}
}

func (r *MapWidgetRenderer) Refresh() {
    r.Layout(r.state.Size())
    canvas.Refresh(r.render)
}

func (r *MapWidgetRenderer) Objects() []fyne.CanvasObject {
    return r.objects
}

func (r *MapWidgetRenderer) Destroy() {
}

func (r *MapWidgetRenderer) draw(w, h int) image.Image {
    im := image.NewRGBA(image.Rect(0, 0, w, h))
    scale := r.state.Scale
    
    for i := 0; i < h; i++ {
        for j := 0; j < w; j++ {
            var c color.Color
            x := 1/scale * float32(j)
            y := 1/scale * float32(i)
            c = r.state.MapViewer.GetPixel(int(x), int(y))
            im.Set(j, i,
                   c)
        }
    }
    r.state.lastImage = im
    return im
}

func (m *MapWidget) CreateRenderer() fyne.WidgetRenderer {
    renderer := &MapWidgetRenderer{state: m}
    
    render := canvas.NewRaster(renderer.draw)
    renderer.render = render
    renderer.objects = []fyne.CanvasObject{render}
    renderer.ApplyTheme()
    
    return renderer
}

func (m *MapWidget) Tapped(ev *fyne.PointEvent) {
    m.clickx = ev.Position.X
    m.clicky = ev.Position.Y
    fmt.Print("x: ", m.clickx, " y:", m.clicky, "\n")
    fmt.Println("size: ", m.size)
    widget.Refresh(m)
}

func (m *MapWidget) TappedSecondary(ev *fyne.PointEvent) {
    widget.Refresh(m)
}
