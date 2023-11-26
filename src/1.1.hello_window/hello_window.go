package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
)

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "LearnOpenGL", nil, nil)
	if err != nil {
		fmt.Println("Failed to create GLFW window: ", err)
		os.Exit(1)
	}
	if window == nil {
		fmt.Println("Failed to create GLFW window: ", err)
		os.Exit(1)
	}

	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width, height int) {
		// TODO: Add code here
		gl.Viewport(0, 0, WindowWidth, WindowHeight)
	})

	err = gl.Init()
	if err != nil {
		fmt.Println("Failed to initialize OpenGL: ", err)
		os.Exit(1)
	}

	for !window.ShouldClose() {
		processInput(window)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	os.Exit(0)
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}
