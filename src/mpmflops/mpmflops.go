package mpmflops

const (
	xval float32 = 0.999950
	aval float32 = 0.000020
	bval float32 = 0.999980
	cval float32 = 0.000011
	dval float32 = 1.000011
	eval float32 = 0.000012
	fval float32 = 0.999992
	gval float32 = 0.000013
	hval float32 = 1.000013
	jval float32 = 0.000014
	kval float32 = 0.999994
	lval float32 = 0.000015
	mval float32 = 1.000015
	oval float32 = 0.000016
	pval float32 = 0.999996
	qval float32 = 0.000017
	rval float32 = 1.000017
	sval float32 = 0.000018
	tval float32 = 1.000018
	uval float32 = 0.000019
	vval float32 = 1.000019
	wval float32 = 0.000021
	yval float32 = 1.000021

	Newdata float32 = 0.999999
)

func TriadPlusLarge(n int, a, b, c, d, e, f, g, h, j, k,
	l, m, o, p, q, r, s, t, u, v, w, y float32, x []float32) {
	for i := 0; i < n; i++ {
		x[i] = (x[i]+a)*b - (x[i]+c)*d + (x[i]+e)*f -
			(x[i]+g)*h + (x[i]+j)*k - (x[i]+l)*m +
			(x[i]+o)*p - (x[i]+q)*r + (x[i]+s)*t -
			(x[i]+u)*v + (x[i]+w)*y
	}
}

func TriadPlusMid(n int, a, b, c, d, e, f float32, x []float32) {
	for i := 0; i < n; i++ {
		x[i] = (x[i]+a)*b - (x[i]+c)*d + (x[i]+e)*f
	}
}

func Triad(n int, a, b float32, x []float32) {
	for i := 0; i < n; i++ {
		x[i] = (x[i] + a) * b
	}
}

func TriadConstPlusLarge(n int, x []float32) {
	for i := 0; i < n; i++ {
		x[i] = (x[i]+aval)*bval - (x[i]+cval)*dval + (x[i]+eval)*fval -
			(x[i]+gval)*hval + (x[i]+jval)*kval - (x[i]+lval)*mval +
			(x[i]+oval)*pval - (x[i]+qval)*rval + (x[i]+sval)*tval -
			(x[i]+uval)*vval + (x[i]+wval)*yval
	}
}

func TriadConstPlusMid(n int, x []float32) {
	for i := 0; i < n; i++ {
		x[i] = (x[i]+aval)*bval - (x[i]+cval)*dval + (x[i]+eval)*fval
	}
}

func TriadConst(n int, x []float32) {
	for i := 0; i < n; i++ {
		x[i] = (x[i] + aval) * bval
	}
}

func getOpwd(part int) int {
	switch part {
	case 0:
		return 2
	case 1:
		return 8
	case 2:
		return 32
	default:
		return 0
	}
}

func initXCpu(value float32, xCpu []float32) {
	for i := 0; i < len(xCpu); i++ {
		xCpu[i] = value
	}
}

func Validate(xCpu []float32, words int, newdata float32) (bool, bool) {
	isok1, isok2 := true, true
	one := xCpu[0]

	if one == newdata {
		isok2, isok1 = false, false
	}

	for i := 1; i < words; i++ {
		if one != xCpu[i] {
			isok1 = false
		}
	}

	return isok1, isok2
}
