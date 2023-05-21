# Optimizing the total stopping time of the Collatz conjecture

Author: Dustin Thériault, 2023

## Abstract

The Collatz conjecture, $3n+1$, is a well-known open problem. Our main result is a function $h$ that empirically reduces the total stopping time of the standard recursive function $f$ and reduced recursive function $g$ for the Collatz conjecture. We have observed the following.

$$
\frac{\sum h(x)}{\sum f(x)} \approx \frac{1}{6} \qquad \frac{\sum h(x)}{\sum g(x)} \approx \frac{1}{2}
$$

We provide source code in Go for all three functions.

## Introduction

We define the standard total stopping time as the following recursive function $f : \mathbb{N}_1 \to \mathbb{N}_0$

$$
f(x) = \begin{cases}
0           & \text{if } x = 1 \\
1 + f(x/2)  & \text{if } x \equiv 0 \pmod 2 \\
1 + f(3x+1) & \text{if } x \equiv 1 \pmod 2 \text{ and } x \ne 1
\end{cases}
$$

This function halts for all $x$ if and only if the Collatz Conjecture is true. If it does halt, then its values are given by the integer sequence [A006577](https://oeis.org/A006577).

The following is a well-known histogram of $f(x)$ up to $10^{10}$

![](results/time_histogram_f_10.png)

```sh
bin/collatz time --graph histogram --fn f --k 10 # command to generate above
```

## Reduced Collatz Function R

The reduced Collatz function $R$ is a well-known technique that simplifies the calculation of the Collatz sequence. The
recursive function $g : \mathbb{N}_1 \to \mathbb{N}_0$ gives the total stopping time of this function, where $2^r$ is
the largest power of 2 that divides the given numerator.

$$
g(x) = \begin{cases}
0                 & \text{if } x = 1 \\
1 + g(x/2^r)      & \text{if } x \equiv 0 \pmod 2 \\
1 + g((3x+1)/2^r) & \text{if } x \equiv 1 \pmod 2 \text{ and } x \ne 1
\end{cases}
$$

If the Collatz conjecture is proven true, then this function halts for all $x$ and its values are given by [A286380](https://oeis.org/A286380).

The following is a histogram of $g(x)$ up to $10^{10}$

![](results/time_histogram_g_10.png)

```sh
bin/collatz time --graph histogram --fn g --k 10 # command to generate above
```

## Main Result

The main result of this paper is the recursive function $h : \mathbb{N}_1 \to \mathbb{N}_0$, which produces a further optimized
total stopping time, where $2^r$ is the largest power of $2$ that divides the given numerator, and $\nu_2(x)$ is the 2-adic valuation.

$$
h(x) = \begin{cases}
0                      & \text{if } x = 1 \\
1 + h(x/2^r)           & \text{if } x \equiv 0 \pmod 2 \\
1 + h((3x+1)/2^r)      & \text{if } x \equiv 1 \pmod 4 \text{ and } x \ne 1 \\
1 + h(k((3x+1)/2)/2^r) & \text{if } x \equiv 3 \pmod 4
\end{cases}
$$

$$
k(x) = (3/2)^{\nu_2(x+1)}(x+1)-1
$$

If the Collatz conjecture is true, then this function halts for all $x$ and its values appear to be given by [A160541](https://oeis.org/A160541).

The following is a histogram of $h(x)$ up to $10^{10}$

![](results/time_histogram_h_10.png)

```sh
bin/collatz time --graph histogram --fn h --k 10 # command to generate above
```

When we combine the histograms of $f(x)$, $g(x)$, and $h(x)$, the improvements become empirically evident.

![](results/time_histogram_combined_10.png)

```sh
bin/collatz time --graph histogram --k 10 # command to generate above
```

A combined scatter plot with $x$ along the x-axis and the values of $f(x), g(x), h(x)$ along the y-axis.

![](results/time_scatter_combined_6.png)

```sh
bin/collatz time --graph scatter --k 6 # command to generate above
```

If we examine only $h(x)$:

![](results/time_scatter_h_6.png)

```sh
bin/collatz time --graph scatter --fn h --k 6 # command to generate above
```

While a rigorous proof of $h$ would be important, we currently do not provide one here. However, one can
recover $f$ using the following variation of $h$, where $2^r$ is the largest factor of 2 that divides a given 
numerator, $k$ is the function defined earlier, and $\nu_2(x)$ is the 2-adic valuation of $x$.

$$
h^\prime(x) = \begin{cases}
0                                                        & \text{if } x = 1 \\
r + h^\prime(x/2^r)                                      & \text{if } x \equiv 0 \pmod 2 \\
1 + r + h^\prime((3x+1)/2^r)                             & \text{if } x \equiv 1 \pmod 4 \text{ and } x \ne 1 \\
2 + r + 2\nu_2(((3x+1)/2)+1) + h^\prime(k((3x+1)/2)/2^r) & \text{if } x \equiv 3 \pmod 4
\end{cases}
$$

$$
h^\prime(x) = f(x)
$$

## Ratios

Based on empirical evidence, the summation of $h(x)$ over the summation of $f(x)$ appears to be around $1/6$.

![](results/ratios_line_f_10.png)

```sh
bin/collatz ratios --graph line --fn f --k 10 --group=1000000 # command to generate above for 10^10 by plotting every millionth ratio
```

![](results/ratios_histogram_f_9.png)

```sh
bin/collatz ratios --graph histogram --fn f --k 9 --group=1000000
```

The ratio between the summation of $h(x)$ over the summation of $g(x)$ appears to approach $1/2$.

![](results/ratios_line_g_10.png)

```sh
bin/collatz ratios --graph line --fn g --k 10 --group=1000000 # command to generate above for 10^10 by plotting every millionth ratio
```

![](results/ratios_histogram_g_9.png)

```sh
bin/collatz ratios --graph histogram --fn g --k 9 --group=10000000
```
