// Package core provides display configuration for the Spectrex framework.
package core

// DisplayConfig holds configuration for window and rendering settings.
type DisplayConfig struct {
	// Window settings
	WindowWidth  int32
	WindowHeight int32
	Title        string
	Maximized    bool
	Resizable    bool
	VSync        bool
	TargetFPS    int32

	// Render settings - if different from window size, rendering is done
	// to a texture and upscaled/downscaled to fit the window
	RenderWidth  int32
	RenderHeight int32

	// Camera defaults
	DefaultFOV float32
}

// DefaultDisplayConfig returns a DisplayConfig with sensible defaults.
func DefaultDisplayConfig() DisplayConfig {
	return DisplayConfig{
		WindowWidth:  1280,
		WindowHeight: 720,
		Title:        "Spectrex",
		Maximized:    false,
		Resizable:    true,
		VSync:        true,
		TargetFPS:    60,
		RenderWidth:  0, // 0 means use window size
		RenderHeight: 0,
		DefaultFOV:   45.0,
	}
}

// EffectiveRenderSize returns the actual render dimensions.
// If RenderWidth/Height are 0, returns window dimensions.
func (dc DisplayConfig) EffectiveRenderSize() (int32, int32) {
	w, h := dc.RenderWidth, dc.RenderHeight
	if w <= 0 {
		w = dc.WindowWidth
	}
	if h <= 0 {
		h = dc.WindowHeight
	}
	return w, h
}
