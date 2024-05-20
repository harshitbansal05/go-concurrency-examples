package mergesort

func Sort(a []int) []int {

	ch := make(chan []int)
	go sort(a, ch)
	return <-ch
}

func sort(a []int, ch chan []int) {
	n := len(a)
	if n == 1 {
		ch <- a
		return
	}

	m := (n - 1) / 2
	if n%2 == 0 {
		m++
	}

	ch1 := make(chan []int)
	ch2 := make(chan []int)
	go sort(a[0:m], ch1)
	go sort(a[m:], ch2)

	var b, c, an []int

	select {
	case b = <-ch1:
		select {
		case c = <-ch2:
		}
	case c = <-ch2:
		select {
		case b = <-ch1:
		}
	}

	i := 0
	j := 0
	for i < m && j < n-m {
		if b[i] < c[j] {
			an = append(an, b[i])
			i++
		} else {
			an = append(an, c[j])
			j++
		}
	}

	for i < m {
		an = append(an, b[i])
		i++
	}

	for j < n-m {
		an = append(an, c[j])
		j++
	}
	ch <- an
}
