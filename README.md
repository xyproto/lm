# linmath

Based on [linmath.h](https://github.com/datenwolf/linmath.h) by Wolfgang 'datenwolf' Draxinger &lt;code@datenwolf.net&gt; (DWTFYWTPL license).

## Features and limitations

* This package is still undergoing heavy refactoring, to make the style more "Go-like".

Several functions and types are provided:

* `Vec2`, a 2 element vector of float64 (x,y)
* `Vec3`, a 3 element vector of float64 (x,y,z)
* `Vec4`, a 4 element vector of float64 (4th component used for homogenous computations)
* `Mat4x4`, a 4 by 4 elements matrix (computations are done in column major order)
* `Quat`, a 4 element vector of float64 (x,y,z,w)

These may be useful for computer graphics programming.

## TODO

- [ ] Add tests

## General info

* License: MIT
* Version: 0.0.0
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
