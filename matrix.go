package lm

import "math"

// Mat4x4 is a [4]Vec4. A [16]float would be faster when using ie. copy.
type Mat4x4 [4]Vec4

// Identity fills the matrix with the identity matrix (diagonal of 1's)
func (M *Mat4x4) Identity() {
	(*M)[0][0] = 1
	(*M)[0][1] = 0
	(*M)[0][2] = 0
	(*M)[0][3] = 0
	(*M)[1][0] = 0
	(*M)[1][1] = 1
	(*M)[1][2] = 0
	(*M)[1][3] = 0
	(*M)[2][0] = 0
	(*M)[2][1] = 0
	(*M)[2][2] = 1
	(*M)[2][3] = 0
	(*M)[3][0] = 0
	(*M)[3][1] = 0
	(*M)[3][2] = 0
	(*M)[3][3] = 1
}

// MatIdentity creates a new identity matrix
func MatIdentity() (r Mat4x4) {
	r.Identity()
	return r
}

// Dup overwrites M with the contents of N
func (M *Mat4x4) Dup(a Mat4x4) {
	copy((*M)[0][:], a[0][:])
	copy((*M)[1][:], a[1][:])
	copy((*M)[2][:], a[2][:])
	copy((*M)[3][:], a[3][:])
}

// Row returns the requested row, counting from 0
func (M *Mat4x4) Row(i int) (r Vec4) {
	for k := 0; k < 4; k++ {
		r[k] = (*M)[k][i]
	}
	return r
}

// Row returns the requested column, counting from 0
func (M *Mat4x4) Col(i int) (r Vec4) {
	copy(r[:], (*M)[i][:])
	return r
}

// Transpose will transpose the given matrix and assign it to the current one
func (M *Mat4x4) Transpose(a Mat4x4) {
	(*M)[0][0] = a[0][0]
	(*M)[1][0] = a[0][1]
	(*M)[2][0] = a[0][2]
	(*M)[3][0] = a[0][3]
	(*M)[0][1] = a[1][0]
	(*M)[1][1] = a[1][1]
	(*M)[2][1] = a[1][2]
	(*M)[3][1] = a[1][3]
	(*M)[0][2] = a[2][0]
	(*M)[1][2] = a[2][1]
	(*M)[2][2] = a[2][2]
	(*M)[3][2] = a[2][3]
	(*M)[0][3] = a[3][0]
	(*M)[1][3] = a[3][1]
	(*M)[2][3] = a[3][2]
	(*M)[3][3] = a[3][3]
}

func (M *Mat4x4) Add(a, b Mat4x4) {
	for i := 0; i < 4; i++ {
		(*M)[i][0] = a[i][0] + b[i][0]
		(*M)[i][1] = a[i][1] + b[i][1]
		(*M)[i][2] = a[i][2] + b[i][2]
		(*M)[i][3] = a[i][3] + b[i][3]
	}
}

func (M *Mat4x4) Sub(a, b Mat4x4) {
	for i := 0; i < 4; i++ {
		(*M)[i][0] = a[i][0] - b[i][0]
		(*M)[i][1] = a[i][1] - b[i][1]
		(*M)[i][2] = a[i][2] - b[i][2]
		(*M)[i][3] = a[i][3] - b[i][3]
	}
}

func (M *Mat4x4) Scale(a Mat4x4, s float64) {
	for i := 0; i < 4; i++ {
		(*M)[i][0] = a[i][0] * s
		(*M)[i][1] = a[i][1] * s
		(*M)[i][2] = a[i][2] * s
		(*M)[i][3] = a[i][3] * s
	}
}

func (M *Mat4x4) ScaleAniso(a Mat4x4, x, y, z float64) {
	(*M)[0] = a[0].Scale(x)
	(*M)[1] = a[1].Scale(y)
	(*M)[2] = a[2].Scale(z)
	(*M)[3][0] = a[3][0]
	(*M)[3][1] = a[3][1]
	(*M)[3][2] = a[3][2]
	(*M)[3][3] = a[3][3]
}

func (M *Mat4x4) Mul(a, b Mat4x4) {
	var temp Mat4x4
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			temp[c][r] = 0
			for k := 0; k < 4; k++ {
				temp[c][r] += a[k][r] * b[c][k]
			}
		}
	}
	M.Dup(temp)
}

func (M *Mat4x4) MulVec4(v Vec4) (r Vec4) {
	for j := 0; j < 4; j++ {
		r[j] = 0
		for i := 0; i < 4; i++ {
			r[j] += (*M)[i][j] * v[i]
		}
	}
	return r
}

func (M *Mat4x4) Translate(x, y, z float64) {
	M.Identity()
	(*M)[3][0] = x
	(*M)[3][1] = y
	(*M)[3][2] = z
}

func (M *Mat4x4) TranslateInPlace(x, y, z float64) {
	t := Vec4{x, y, z, 0}
	for i := 0; i < 4; i++ {
		(*M)[3][i] += M.Row(i).MulInner(t)
	}
}

func (M *Mat4x4) FromVec3MulOuter(a, b Vec3) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if i < 3 && j < 3 {
				(*M)[i][j] = a[i] * b[j]
			} else {
				(*M)[i][j] = 0
			}
		}
	}
}

func (M *Mat4x4) Rotate(a Mat4x4, x, y, z, angle float64) {
	s := math.Sin(angle)
	c := math.Cos(angle)
	u := Vec3{x, y, z}
	T := &Mat4x4{}

	if u.Len() > 1e-4 {
		u = u.Norm()
		T.FromVec3MulOuter(u, u)

		S := &Mat4x4{
			Vec4{0, u[2], -u[1], 0},
			Vec4{-u[2], 0, u[0], 0},
			Vec4{u[1], -u[0], 0, 0},
			Vec4{0, 0, 0, 0},
		}
		S.Scale(*S, s)

		C := &Mat4x4{}
		C.Identity()
		C.Sub(*C, *T)
		C.Scale(*C, c)

		T.Add(*T, *C)
		T.Add(*T, *S)

		T[3][3] = 1.0
		M.Mul(a, *T)
	} else {
		M.Dup(a)
	}
}

func (M *Mat4x4) RotateX(a Mat4x4, angle float64) {
	s := math.Sin(angle)
	c := math.Cos(angle)
	b := Mat4x4{
		Vec4{1.0, 0, 0, 0},
		Vec4{0, c, s, 0},
		Vec4{0, -s, c, 0},
		Vec4{0, 0, 0, 1.0},
	}
	M.Mul(a, b)
}

func (M *Mat4x4) RotateY(a Mat4x4, angle float64) {
	s := math.Sin(angle)
	c := math.Cos(angle)
	b := Mat4x4{
		Vec4{c, 0, s, 0},
		Vec4{0, 1.0, 0, 0},
		Vec4{-s, 0, c, 0},
		Vec4{0, 0, 0, 1.0},
	}
	M.Mul(a, b)
}

func (M *Mat4x4) RotateZ(a Mat4x4, angle float64) {
	s := math.Sin(angle)
	c := math.Cos(angle)
	b := Mat4x4{
		Vec4{c, s, 0, 0},
		Vec4{-s, c, 0, 0},
		Vec4{0, 0, 1.0, 0},
		Vec4{0, 0, 0, 1.0},
	}
	M.Mul(a, b)
}

func (M *Mat4x4) Invert(a Mat4x4) {
	var s, c [6]float64

	s[0] = a[0][0]*a[1][1] - a[1][0]*a[0][1]
	s[1] = a[0][0]*a[1][2] - a[1][0]*a[0][2]
	s[2] = a[0][0]*a[1][3] - a[1][0]*a[0][3]
	s[3] = a[0][1]*a[1][2] - a[1][1]*a[0][2]
	s[4] = a[0][1]*a[1][3] - a[1][1]*a[0][3]
	s[5] = a[0][2]*a[1][3] - a[1][2]*a[0][3]

	c[0] = a[2][0]*a[3][1] - a[3][0]*a[2][1]
	c[1] = a[2][0]*a[3][2] - a[3][0]*a[2][2]
	c[2] = a[2][0]*a[3][3] - a[3][0]*a[2][3]
	c[3] = a[2][1]*a[3][2] - a[3][1]*a[2][2]
	c[4] = a[2][1]*a[3][3] - a[3][1]*a[2][3]
	c[5] = a[2][2]*a[3][3] - a[3][2]*a[2][3]

	// Assumes it is invertible
	idet := 1.0 / (s[0]*c[5] - s[1]*c[4] + s[2]*c[3] + s[3]*c[2] - s[4]*c[1] + s[5]*c[0])

	(*M)[0][0] = (a[1][1]*c[5] - a[1][2]*c[4] + a[1][3]*c[3]) * idet
	(*M)[0][1] = (-a[0][1]*c[5] + a[0][2]*c[4] - a[0][3]*c[3]) * idet
	(*M)[0][2] = (a[3][1]*s[5] - a[3][2]*s[4] + a[3][3]*s[3]) * idet
	(*M)[0][3] = (-a[2][1]*s[5] + a[2][2]*s[4] - a[2][3]*s[3]) * idet

	(*M)[1][0] = (-a[1][0]*c[5] + a[1][2]*c[2] - a[1][3]*c[1]) * idet
	(*M)[1][1] = (a[0][0]*c[5] - a[0][2]*c[2] + a[0][3]*c[1]) * idet
	(*M)[1][2] = (-a[3][0]*s[5] + a[3][2]*s[2] - a[3][3]*s[1]) * idet
	(*M)[1][3] = (a[2][0]*s[5] - a[2][2]*s[2] + a[2][3]*s[1]) * idet

	(*M)[2][0] = (a[1][0]*c[4] - a[1][1]*c[2] + a[1][3]*c[0]) * idet
	(*M)[2][1] = (-a[0][0]*c[4] + a[0][1]*c[2] - a[0][3]*c[0]) * idet
	(*M)[2][2] = (a[3][0]*s[4] - a[3][1]*s[2] + a[3][3]*s[0]) * idet
	(*M)[2][3] = (-a[2][0]*s[4] + a[2][1]*s[2] - a[2][3]*s[0]) * idet

	(*M)[3][0] = (-a[1][0]*c[3] + a[1][1]*c[1] - a[1][2]*c[0]) * idet
	(*M)[3][1] = (a[0][0]*c[3] - a[0][1]*c[1] + a[0][2]*c[0]) * idet
	(*M)[3][2] = (-a[3][0]*s[3] + a[3][1]*s[1] - a[3][2]*s[0]) * idet
	(*M)[3][3] = (a[2][0]*s[3] - a[2][1]*s[1] + a[2][2]*s[0]) * idet
}

func (M *Mat4x4) Orthonormalize(a Mat4x4) {
	M.Dup(a)

	vn := (*M)[2].Vec3().Norm()

	(*M)[2][0] = vn[0]
	(*M)[2][1] = vn[1]
	(*M)[2][2] = vn[2]

	s := (*M)[1].Vec3().MulInner((*M)[2].Vec3())
	h := (*M)[2].Vec3().Scale(s)
	vs := (*M)[1].Vec3().Sub(h)

	(*M)[1][0] = vs[0]
	(*M)[1][1] = vs[1]
	(*M)[1][2] = vs[2]

	vn = (*M)[1].Vec3().Norm()

	(*M)[1][0] = vn[0]
	(*M)[1][1] = vn[1]
	(*M)[1][2] = vn[2]

	s = (*M)[0].Vec3().MulInner((*M)[2].Vec3())
	h = (*M)[2].Vec3().Scale(s)
	vs = (*M)[0].Vec3().Sub(h)

	(*M)[0][0] = vs[0]
	(*M)[0][1] = vs[1]
	(*M)[0][2] = vs[2]

	s = (*M)[0].Vec3().MulInner((*M)[1].Vec3())
	h = (*M)[1].Vec3().Scale(s)
	vn = (*M)[0].Vec3().Sub(h).Norm()

	(*M)[0][0] = vn[0]
	(*M)[0][1] = vn[1]
	(*M)[0][2] = vn[2]
}

func (M *Mat4x4) Frustum(l, r, b, t, n, f float64) {
	(*M)[0][0] = 2.0 * n / (r - l)
	(*M)[0][1] = 0
	(*M)[0][2] = 0
	(*M)[0][3] = 0

	(*M)[1][1] = 2. * n / (t - b)
	(*M)[1][0] = 0
	(*M)[1][2] = 0
	(*M)[1][3] = 0

	(*M)[2][0] = (r + l) / (r - l)
	(*M)[2][1] = (t + b) / (t - b)
	(*M)[2][2] = -(f + n) / (f - n)
	(*M)[2][3] = -1.0

	(*M)[3][2] = -2.0 * (f * n) / (f - n)
	(*M)[3][0] = 0
	(*M)[3][1] = 0
	(*M)[3][3] = 0
}

func (M *Mat4x4) Ortho(l, r, b, t, n, f float64) {
	(*M)[0][0] = 2.0 / (r - l)
	(*M)[0][1] = 0
	(*M)[0][2] = 0
	(*M)[0][3] = 0

	(*M)[1][1] = 2.0 / (t - b)
	(*M)[1][0] = 0
	(*M)[1][2] = 0
	(*M)[1][3] = 0

	(*M)[2][2] = -2.0 / (f - n)
	(*M)[2][0] = 0
	(*M)[2][1] = 0
	(*M)[2][3] = 0

	(*M)[3][0] = -(r + l) / (r - l)
	(*M)[3][1] = -(t + b) / (t - b)
	(*M)[3][2] = -(f + n) / (f - n)
	(*M)[3][3] = 1.0
}

func (M *Mat4x4) Perspective(yFOV, aspect, n, f float64) {
	/* NOTE: Degrees are an unhandy unit to work with.
	 * linmath.h uses radians for everything! */
	a := 1.0 / math.Tan(yFOV/2.0)

	(*M)[0][0] = a / aspect
	(*M)[0][1] = 0
	(*M)[0][2] = 0
	(*M)[0][3] = 0

	(*M)[1][0] = 0
	(*M)[1][1] = a
	(*M)[1][2] = 0
	(*M)[1][3] = 0

	(*M)[2][0] = 0
	(*M)[2][1] = 0
	(*M)[2][2] = -((f + n) / (f - n))
	(*M)[2][3] = -1.0

	(*M)[3][0] = 0
	(*M)[3][1] = 0
	(*M)[3][2] = -((2.0 * f * n) / (f - n))
	(*M)[3][3] = 0
}

func (M *Mat4x4) LookAt(eye, center, up Vec3) {
	/* Adapted from Android's OpenGL Matrix.java.                        */
	/* See the OpenGL GLUT documentation for gluLookAt for a description */
	/* of the algorithm. We implement it in a straightforward way:       */

	/* TODO: The negation of of can be spared by swapping the order of
	 *       operands in the following cross products in the right way. */

	f := center.Sub(eye).Norm()
	s := f.MulCross(up).Norm()
	t := s.MulCross(f)

	(*M)[0][0] = s[0]
	(*M)[0][1] = t[0]
	(*M)[0][2] = -f[0]
	(*M)[0][3] = 0

	(*M)[1][0] = s[1]
	(*M)[1][1] = t[1]
	(*M)[1][2] = -f[1]
	(*M)[1][3] = 0

	(*M)[2][0] = s[2]
	(*M)[2][1] = t[2]
	(*M)[2][2] = -f[2]
	(*M)[2][3] = 0

	(*M)[3][0] = 0
	(*M)[3][1] = 0
	(*M)[3][2] = 0
	(*M)[3][3] = 1.0

	M.TranslateInPlace(-eye[0], -eye[1], -eye[2])
}
