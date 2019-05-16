/**
  * Author: JeffreyBool
  * Date: 2019/5/16
  * Time: 17:10
  * Software: GoLand
*/

package core

import (
	"fmt"
)

/**
 * AOI 区域管理模块
 */
type AOIManager struct {
	MinX  int           //区域左边界坐标
	MaxX  int           //区域右边界坐标
	MinY  int           //区域上边界坐标
	MaxY  int           //区域下边界坐标
	CntsX int           //x方向格子的数量
	CntsY int           //y方向的格子数量
	Grids map[int]*Grid //当前区域中都有哪些格子，key=格子ID， value=格子对象
}

// 初始化一个AOI区域
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		MinY:  minY,
		MaxY:  maxY,
		CntsX: cntsX,
		CntsY: cntsY,
		Grids: make(map[int]*Grid),
	}

	//给AOI初始化区域中所有的格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			//计算格子ID
			//格子编号：id = idy *nx + idx  (利用格子坐标得到格子编号)
			gid := y*cntsX + x

			//初始化一个格子放在AOI中的map里，key是当前格子的ID
			aoiMgr.Grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}

	return aoiMgr
}

//得到每个格子在 X 轴的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

//得到每个格子在 Y 轴的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

//根据格子的gID得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGid(gid int) (grids []*Grid) {
	//判断gID是否存在
	if _, ok := m.Grids[gid]; !ok {
		return
	}

	//将当前gid添加到九宫格中
	grids = append(grids, m.Grids[gid])

	//根据gid得到当前格子所在的X轴编号(idx = id % nx)
	idx := gid % m.CntsX

	//判断当前idx左边是否还有格子
	if idx > 0 {
		grids = append(grids, m.Grids[gid-1])
	}

	//判断当前的idx右边是否还有格子
	if idx < m.CntsX-1 {
		grids = append(grids, m.Grids[gid+1])
	}

	//将x轴当前的格子都取出，进行遍历，再分别得到每个格子的上下是否有格子

	//得到当前x轴的格子id集合
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.Id)
	}

	//遍历x轴格子
	for _, v := range gidsX {
		//计算该格子处于第几列
		idy := v / m.CntsX

		//判断当前的idy上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.Grids[v-m.CntsX])
		}
		//判断当前的idy下边是否还有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.Grids[v+m.CntsX])
		}
	}
	return
}

//通过横纵坐标得到周边九宫格内的全部PlayerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIds []int) {
	//根据横纵坐标得到当前坐标属于哪个格子ID
	gID := m.GetGidByPos(x, y)

	//根据格子ID得到周边九宫格的信息
	grids := m.GetSurroundGridsByGid(gID)
	for _, v := range grids {
		playerIds = append(playerIds, v.GetPlayerIds()...)
	}

	return
}

//通过横纵坐标获取对应的格子ID
func (m *AOIManager) GetGidByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(x) - m.MinY) / m.gridLength()

	return gy*m.CntsX + gx
}

//通过GID获取当前格子的全部playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.Grids[gID].GetPlayerIds()
	return
}

//移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.Grids[gID].Remove(pID)
}

//添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.Grids[gID].Add(pID)
}

//通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.Grids[gID]
	grid.Add(pID)
}

//通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.Grids[gID]
	grid.Remove(pID)
}

//打印信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\n minX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d \n Grids in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.Grids {
		s += fmt.Sprintln(grid)
	}

	return s
}
