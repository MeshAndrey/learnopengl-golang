package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
)

const vertexShaderSource string = `
#version 330 core
layout (location = 0) in vec3 aPos;
void main()
{
	gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
};
`

const fragmentShaderSource string = `
#version 330 core
out vec4 FragColor;
void main()
{
    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
};
`

func init() {
	// GLFW event handling must be run on the main OS thread
	runtime.LockOSThread()
}

func checkShaderCompileErrors(shader uint32) {
	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		var infoLog [512]byte
		gl.GetShaderInfoLog(shader, 512, nil, (*uint8)(unsafe.Pointer(&infoLog)))
		fmt.Println("ERROR::SHADER::COMPILATION_FAILED\n", string(infoLog[:512]))
	}
}

func checkProgramLinkErrors(program uint32) {
	var success int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &success)
	if success == 0 {
		var infoLog [512]byte
		gl.GetProgramInfoLog(program, 512, nil, (*uint8)(unsafe.Pointer(&infoLog)))
		fmt.Println("ERROR::PROGRAM::LINK_FAILED\n", string(infoLog[:512]))
	}
}

func compileShader() []uint32 {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	shaderSourceChars, freeVertexShaderFunc := gl.Strs(vertexShaderSource)
	defer freeVertexShaderFunc()
	gl.ShaderSource(vertexShader, 1, shaderSourceChars, nil)
	gl.CompileShader(vertexShader)
	checkShaderCompileErrors(vertexShader)

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	shaderSourceChars, freeFragmentShaderFunc := gl.Strs(fragmentShaderSource)
	defer freeFragmentShaderFunc()
	gl.ShaderSource(fragmentShader, 1, shaderSourceChars, nil)
	gl.CompileShader(fragmentShader)
	checkShaderCompileErrors(fragmentShader)

	return []uint32{vertexShader, fragmentShader}
}

func linkShaders(shaders []uint32) uint32 {
	program := gl.CreateProgram()
	for _, shader := range shaders {
		gl.AttachShader(program, shader)
	}
	gl.LinkProgram(program)
	checkProgramLinkErrors(program)

	// shader objects are not needed after they are linked into a program object
	for _, shader := range shaders {
		gl.DeleteShader(shader)
	}

	return program
}

func createTriangleVAO() uint32 {
	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// specify the format of our vertex input
	// (shader) input 0
	// vertex has size 3
	// vertex items are of type FLOAT
	// do not normalize (already done)
	// stride of 3 * sizeof(float) (separation of vertices)
	// offset of where the position data starts (0 for the beginning)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	return VAO
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

	shaders := compileShader()
	shaderProgram := linkShaders(shaders)

	VAO := createTriangleVAO()

	for !window.ShouldClose() {
		processInput(window)

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shaderProgram)
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	gl.DeleteVertexArrays(1, &VAO)
	// gl.DeleteBuffers(1, &VBO)
	gl.DeleteProgram(shaderProgram)

	os.Exit(0)
}

func processInput(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}
