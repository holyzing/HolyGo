import numpy as np
t1= np.arange(12).reshape(3,4).astype(float)

t1[1, 2:]=np.nan

print(t1)
print(t1.shape[1])
# 数组的列数，t1.shape[0]行数

def fill_ndarray(t):
    for i in range(t.shape[1]): #遍历每一列
        temp_col = t[:,i] # 当前列
        nan_num = np.count_nonzero(temp_col!=temp_col)


        print(nan_num)
        if nan_num != 0: # 不为零，说明当前这一列中有nan

            # np.nan != np.nan
            temp_not_nan_col = temp_col[temp_col==temp_col]

            # 当前一列不为nan的array
            # 选中当前为nan的位置，把值赋值为不为nan的均值
            temp_col[np.isnan(temp_col)] = temp_not_nan_col.mean()
            t[:, i] = temp_col
    return t


if __name__ == '__main__':
    t2=np.arange(24).reshape(4,6).astype(float)
    t2[1,2:]=np.nan
    print(t2)

    print('*'*100)
    fill_ndarray(t2)
    print(t2)
