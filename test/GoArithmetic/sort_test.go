package GoArithmetic

import (
	"fmt"
	"testing"
)

func TestQuickSort(t *testing.T) {
	/*
		遍历一次的时间复杂度：O(N)
		遍历次数：最多 N 次，最少 lg(N+1) 次
		平均时间复杂度： O(N*lg(N))
		最坏时间复杂度： O(N^2)

		算法稳定性：
			假设在数列中存在a[i]=a[j]，若在排序之前，a[i] 在 a[j] 前面，
			并且在排序之后，a[i] 仍然在 a[j] 前面，则这个排序算法才算是稳定的。

		二叉树的 深度K 和 节点数N 的关系：2^k-1 <= N

		分冶法：
			把一个复杂的问题分解成两个或更多的相同或相似的子问题，再把子问题分成更小的子问题，
			直到最后子问题可以简单的直接求解 （各个击破），原问题的解即子问题的解的合并。
			第一步（分）：将原来复杂的问题分解为若干个规模较小、相互独立、
						与原问题形式相同的子问题，分解到可以直接求解为止；
			第二步（治）：此时可以直接求解；
			第三步（合）：将小规模的问题的解合并为一个更大规模的问题的解，自底向上逐步求出原来问题的解
	*/

	Arr := []int{23, 65, 13, 27, 42, 15, 38, 21, 4, 10}
	qsort(Arr, 0, len(Arr)-1)
	fmt.Println(Arr)
}

/**
快速排序：分治法+递归实现随意取一个值A，将比A大的放在A的右边，比A小的放在A的左边；
然后在左边的值AA中再取一个值B，将AA中比B小的值放在B的左边，将比B大的值放在B的右边。以此类推*/
func qsort(arr []int, first, last int) {
	flag := first
	left := first
	right := last
	if first >= last {
		return
	} // 将大于arr[flag]的都放在右边，小于的，都放在左边
	for first < last {
		// 如果flag从左边开始，那么是必须先从有右边开始比较，也就是先在右边找比flag小的
		for first < last {
			if arr[last] >= arr[flag] {
				last--
				continue
			}
			// 交换数据
			arr[last], arr[flag] = arr[flag], arr[last]
			flag = last
			break
		}
		for first < last {
			if arr[first] <= arr[flag] {
				first++
				continue
			}
			arr[first], arr[flag] = arr[flag], arr[first]
			flag = first
			break
		}
	}
	qsort(arr, left, flag-1)
	qsort(arr, flag+1, right)
}
