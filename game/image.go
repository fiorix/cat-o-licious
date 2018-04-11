// Copyright 2017  The cat-o-licious authors.
//
// Licensed under GNU General Public License 3.0.
// Some rights reserved. See LICENSE, AUTHORS.

package game

import (
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
	sdlimg "github.com/veandco/go-sdl2/img"
)

// Image is a wrapper for SDL images.
type Image interface {
	// Size returns the image size.
	Size() (w, h int32)

	// Texture returns the SDL image texture.
	Texture() *sdl.Texture
}

type image struct {
	w, h int32
	t    *sdl.Texture
}

// NewImageFromFile loads an image from file.
func NewImageFromFile(r *sdl.Renderer, file string) (Image, error) {
	img, err := sdlimg.LoadTexture(r, file)
	if err != nil {
		return nil, err
	}
	_, _, w, h, err := img.Query()
	if err != nil {
		return nil, err
	}
	return &image{
		w: w,
		h: h,
		t: img,
	}, nil
}

// Size implements the Image interface.
func (img *image) Size() (w, h int32) {
	return img.w, img.h
}

// Texture returns the SDL image's texture.
func (img *image) Texture() *sdl.Texture {
	return img.t
}

// NewImageSetFromFiles loads a sequence of images with the same prefix.
// The format is {prefix}{n}.png where N starts from 1. At least one
// image must exist otherwise an error is returned.
func NewImageSetFromFiles(r *sdl.Renderer, prefix string) ([]Image, error) {
	i := 0
	var imgs []Image
	for {
		idx := strconv.Itoa(i + 1)
		img, err := NewImageFromFile(r, prefix+idx+".png")
		if err != nil {
			if i == 0 {
				return nil, err
			}
			break
		}
		imgs = append(imgs, img)
		i++
	}
	return imgs, nil
}
