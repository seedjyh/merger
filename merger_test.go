// merger_test
package merger

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Payment struct {
	id         string    // 唯一性标记，不参与分组
	categoryID int       // 类型序号
	account    string    // 必须完全相同
	payTime    time.Time // 需要相近
}

func (p *Payment) CategoryID() int {
	return p.categoryID
}

func (p *Payment) Key() string {
	return p.account
}

func (p *Payment) DifferenceFrom(other Item) float64 {
	return float64(p.payTime.Sub(other.(*Payment).payTime))
}

func TestMerger_Merge(t *testing.T) {
	// 没有加入任何Item，则结果是空的
	m := NewMerger(3)
	assert.Equal(t, 0, len(m.Merge(1.0)))
}

func TestMerger_Merge1(t *testing.T) {
	// 加入1个Item，则结果就是那个
	m := NewMerger(3)
	testData := []*Payment{{
		id:         "a",
		categoryID: 1,
		account:    "jyh",
		payTime:    time.Now(),
	}}
	for _, item := range testData {
		assert.NoError(t, m.AppendItem(item))
	}
	groups := m.Merge(1.0)
	assert.Equal(t, 1, len(groups))
	assert.Equal(t, nil, groups[0].Items[0])
	assert.Equal(t, testData[0], groups[0].Items[1])
	assert.Equal(t, nil, groups[0].Items[2])
}

func TestMerger_Merge2(t *testing.T) {
	// 加入2个Item，相同的key，不同的category，但模糊匹配指数差异较大，则会出现在不同的group中
	m := NewMerger(3)
	testData := []*Payment{{
		id:         "a",
		categoryID: 1,
		account:    "jyh",
		payTime:    time.Now(),
	},{
		id:         "b",
		categoryID: 2,
		account:    "jyh",
		payTime:    time.Now().Add(time.Hour),
	}}
	for _, item := range testData {
		assert.NoError(t, m.AppendItem(item))
	}
	groups := m.Merge(1.0)
	assert.Equal(t, 2, len(groups))
	// groups[0]
	assert.Equal(t, nil, groups[0].Items[0])
	assert.Equal(t, testData[0], groups[0].Items[1])
	assert.Equal(t, nil, groups[0].Items[2])
	// groups[1]
	assert.Equal(t, nil, groups[1].Items[0])
	assert.Equal(t, nil, groups[1].Items[1])
	assert.Equal(t, testData[1], groups[1].Items[2])
}

func TestMerger_Merge3(t *testing.T) {
	// 加入2个Item，相同的key，不同的category，但模糊匹配指数差异较小，则会出现在同一个group中
	m := NewMerger(3)
	testData := []*Payment{{
		id:         "a",
		categoryID: 1,
		account:    "jyh",
		payTime:    time.Now(),
	},{
		id:         "b",
		categoryID: 2,
		account:    "jyh",
		payTime:    time.Now().Add(time.Nanosecond),
	}}
	for _, item := range testData {
		assert.NoError(t, m.AppendItem(item))
	}
	groups := m.Merge(float64(time.Hour))
	assert.Equal(t, 1, len(groups))
	// groups[0]
	assert.Equal(t, nil, groups[0].Items[0])
	assert.Equal(t, testData[0], groups[0].Items[1])
	assert.Equal(t, testData[1], groups[0].Items[2])
}

func TestMerger_MergeX(t *testing.T) {
	// 纵轴是来源，横轴是时间。所有item的key相同。
	//	\	0	1	2	3	4	5
	//	0	0a	1b	2c	3d	.	.
	//	1	.	4e	.	5f	6g	.
	//	2	.	7h	8i	.	.	9j
	m := NewMerger(3)
	testData := make([]*Payment, 10)
	for i := 0; i < 10; i++ {
		testData[i] = &Payment{
			id:         string(rune('a' + i)),
			account:    "jyh",
		}
	}
	startTime := time.Now()
	testData[0].categoryID = 0
	testData[1].categoryID = 0
	testData[2].categoryID = 0
	testData[3].categoryID = 0
	testData[4].categoryID = 1
	testData[5].categoryID = 1
	testData[6].categoryID = 1
	testData[7].categoryID = 2
	testData[8].categoryID = 2
	testData[9].categoryID = 2
	testData[0].payTime = startTime
	testData[1].payTime = startTime.Add(time.Minute + time.Nanosecond)
	testData[4].payTime = startTime.Add(time.Minute + time.Millisecond)
	testData[7].payTime = startTime.Add(time.Minute + time.Second)
	testData[2].payTime = startTime.Add(time.Minute * 2 + time.Nanosecond)
	testData[8].payTime = startTime.Add(time.Minute * 2 + time.Millisecond)
	testData[3].payTime = startTime.Add(time.Minute * 3 + time.Second)
	testData[5].payTime = startTime.Add(time.Minute * 3 + time.Nanosecond)
	testData[6].payTime = startTime.Add(time.Minute * 4 + time.Millisecond)
	testData[9].payTime = startTime.Add(time.Minute * 5 + time.Second)
	for _, p := range testData {
		assert.NoError(t, m.AppendItem(p))
	}
	groups := m.Merge(float64(time.Second * 10))
	assert.Equal(t, 6, len(groups))
	var nowGroup *ItemGroup
	nowGroup = groups[0]
	assert.Equal(t, testData[0], nowGroup.Items[0])
	assert.Equal(t, nil, nowGroup.Items[1])
	assert.Equal(t, nil, nowGroup.Items[2])
	nowGroup = groups[1]
	assert.Equal(t, testData[1], nowGroup.Items[0])
	assert.Equal(t, testData[4], nowGroup.Items[1])
	assert.Equal(t, testData[7], nowGroup.Items[2])
	nowGroup = groups[2]
	assert.Equal(t, testData[2], nowGroup.Items[0])
	assert.Equal(t, nil, nowGroup.Items[1])
	assert.Equal(t, testData[8], nowGroup.Items[2])
	nowGroup = groups[3]
	assert.Equal(t, testData[3], nowGroup.Items[0])
	assert.Equal(t, testData[5], nowGroup.Items[1])
	assert.Equal(t, nil, nowGroup.Items[2])
	nowGroup = groups[4]
	assert.Equal(t, nil, nowGroup.Items[0])
	assert.Equal(t, testData[6], nowGroup.Items[1])
	assert.Equal(t, nil, nowGroup.Items[2])
	nowGroup = groups[5]
	assert.Equal(t, nil, nowGroup.Items[0])
	assert.Equal(t, nil, nowGroup.Items[1])
	assert.Equal(t, testData[9], nowGroup.Items[2])
}
