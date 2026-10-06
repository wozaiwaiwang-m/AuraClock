package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"magpie"
	"os"
	"strconv"
	"text/tabwriter"

	// "theme"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	// "fyne.io/fyne/v2/driver/desktop"
	"github.com/gen2brain/beeep"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"

	// "github.com/hajimehoshi/oto/v3"

	// "fyne.io/fyne/v2/dialog"
	// "fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/bradfitz/latlong"
	"github.com/nathan-osman/go-sunrise"
)

type AuraClock struct {
	Name    string        `json:"name"`    // 闹钟名称
	Sun     string        `json:"sun"`     // 日出（sunRise）或日落（sunSet）
	Diff    time.Duration `json:"diff"`    // 时间差值（分钟）
	Remark  string        `json:"remark"`  // 闹钟的备注
	Enabled bool          `json:"enabled"` // 闹钟的开关
}

var (
	ReadPath  string
	WritePath string

	onlyCreate bool
	guiMode    bool
	remind     bool

	auraList  []AuraClock
	longitude float64 // 经度
	latitude  float64 // 纬度
	sunRise   time.Time
	sunSet    time.Time
	// nextClock time.Duration //现在与闹钟时间的差值
)

func StartAlarmRoutines(auraList []AuraClock) {
	var nextClock time.Time
	var diffClock time.Duration
	// var aura AuraClock
	for _, aura := range auraList {

		nextClock = ComClock(aura)
		diffClock = nextClock.Sub(time.Now())

		if nextClock.Before(time.Now()) {
			diffClock = 24*time.Hour - diffClock
		}
		go func() {
			time.Sleep(diffClock)
			TriggerClock(aura)

		}()
		// switch:
		// case aura.

	}
}

func TriggerClock(aura AuraClock) {
	err := beeep.Notify("闹钟提醒："+aura.Name, aura.Remark, "aads")
	if err != nil {
		fmt.Println("发送通知失败：", err)
	}
	//声音提醒
	audioPath := "data/oga/闹钟.mp3"

	audioFile, err := os.Open(audioPath)
	if err != nil {
		log.Fatal(err)
	}
	defer audioFile.Close()

	mp3Decoder, err := mp3.NewDecoder(audioFile)
	if err != nil {
		log.Fatal(err)
	}

	ctx, ready, err := oto.NewContext(
		mp3Decoder.SampleRate(),
		2,
		oto.FormatSignedInt16LE,
	)
	if err != nil {
		log.Fatal(err)
	}

	<-ready

	player := ctx.NewPlayer(mp3Decoder)
	defer player.Close()

	player.Play()

	for player.IsPlaying() {
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("已发送提醒：", aura.Name)
}

func fmtJson(clocks []AuraClock) (string, []byte) {
	data, err := json.MarshalIndent(clocks, "", "	") // 序列化
	if err != nil {
		fmt.Println("序列化失败：", err)
	}
	return string(data), data
}

// func UpdateFile(jsonByte []byte, file string) {
func UpdateFile(alist []AuraClock, file string) {
	var list []AuraClock
	for _, aura := range alist {
		if aura.Name != "" {
			list = append(list, aura)
		}
	}
	// var outlist []AuraClock
	err := os.Truncate(file, 0)
	if err != nil {
		fmt.Println("清空失败")
	}

	WriteConfig(list, file)
	// os.WriteFile(file, jsonByte, 0o644)
}

func SortAuraList(list *[]AuraClock) { // 通用的按字母排序的函数
	sort.Slice(*list, func(i, j int) bool {
		return (*list)[i].Name < (*list)[j].Name
	})
}

func StartGUI() {
	os.Setenv("FYNE_FONT", "/usr/share/fonts/google-droid-sans-fonts/DroidSansFallbackFull.ttf") // 设置字体，防止乱码

	sunMap := map[string]string{
		"日出": "sunRise",
		"日落": "sunSet",
	}

	// 新建窗口
	auraGUI := app.New()
	auraWindow := auraGUI.NewWindow("光律生活")

	var (
		width  float32 = 600
		height float32 = 400
	)

	auraWindow.Resize(fyne.NewSize(width, height))

	auraList = ReadConfig("GUI.json")

	// auraList = <-red

	// 列表组件
	list := widget.NewList(
		// 根据已有的闹钟列表定义一共有多少行
		func() int {
			return len(auraList)
		},

		func() fyne.CanvasObject { // 定义每一行有两列
			return container.NewGridWithColumns(3,
				widget.NewLabel(""),
				widget.NewLabel(""),
				widget.NewButton("", nil),
			)
		},

		func(lii widget.ListItemID, co fyne.CanvasObject) {
			labels := co.(*fyne.Container).Objects // 一个列表
			nameLabel := labels[0].(*widget.Label)
			timeLabel := labels[1].(*widget.Label)
			deleteBtn := labels[2].(*widget.Button)
			// aura := auraList

			deleteBtn.OnTapped = func() {
				// 把已经删除的行涂灰
				nameLabel.SetText("已移除")
				timeLabel.SetText("已移除")
				deleteBtn.Disable()
				deleteBtn.SetText("已移除")

				auraList[lii] = AuraClock{}
			}
			// aura = append([]AuraClock{aura[lii]}, auraList...)
			nameLabel.SetText(auraList[lii].Name)
			timeLabel.SetText(FormatClock(auraList[lii]))
		},
	)

	// 输入组件
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("闹钟名称...")

	diffEntry := widget.NewEntry()
	diffEntry.SetPlaceHolder("偏移分钟 (如 -30 或 10)...")

	var selectedSun string = "sunRise"
	selectBox := widget.NewSelect([]string{"日出", "日落"}, func(s string) {
		selectedSun = sunMap[s]
	})
	selectBox.PlaceHolder = "请选择日出或日落"
	// selectBox.SetSelected("sunRise")
	remarkEntry := widget.NewEntry()
	remarkEntry.SetPlaceHolder("闹钟备注……")

	// 提交按钮逻辑
	referButton := widget.NewButton("添加闹钟", func() {
		diffInt, _ := strconv.Atoi(diffEntry.Text)
		newAura := AuraClock{
			Name:   nameEntry.Text,
			Sun:    selectedSun,
			Diff:   time.Duration(diffInt) * time.Minute,
			Remark: remarkEntry.Text,
		}

		// out <- newAura

		// 为了简单起见，这里直接追加并刷新，也可以重新 ReadConfig
		if newAura.Name != "" {
			auraList = append([]AuraClock{newAura}, auraList...)
			SortAuraList(&auraList)
		} else {
			// dialog.ShowInformation("错误", "空闹钟", auraWindow)
			popup := widget.NewPopUp(widget.NewLabel("空闹钟"), auraWindow.Canvas())
			popup.Move(fyne.NewPos(width/3, height/3))
			popup.Show()
			fmt.Println("空闹钟，不保存")
		}

		// UpdateFile(auraList, "GUI.json")
		list.Refresh()

		nameEntry.SetText("") // 清空名称输入框
		diffEntry.SetText("") // 清空差值输入框
		remarkEntry.SetText("")

		// nameEntry.SetText("")
		// diffEntry.SetText("")
	})

	// --- 布局设置 ---
	// 左侧输入区
	inputBox := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("新增闹钟"),
			nameEntry,
			diffEntry,
			selectBox),

		referButton,
		nil,
		nil,
		remarkEntry,
	)

	// 整体布局：左边输入，右边列表（带滚动条）
	// 用 Scroll 容器包裹列表防止闹钟多了以后溢出
	content := container.NewHSplit(
		inputBox,
		container.NewScroll(list),
	)
	content.Offset = 0.4 // 设置左右比例

	auraWindow.SetContent(content)

	StartAlarmRoutines(auraList)

	auraWindow.ShowAndRun()
}

func CreateClock() AuraClock { // 创建闹钟，并输出一个AuraClock结构体
	// 在命令行中获取用户输入
	na := magpie.AskString()
	_, su := magpie.AskArray([]string{"sunRise", "sunSet"})
	mi := magpie.AskInt()
	di := time.Duration(mi) * time.Minute
	// 返回新建的闹钟
	return AuraClock{Name: na, Sun: su, Diff: di}
}

func ReadConfig(file string) []AuraClock { // 读取文件并输出
	var list []AuraClock

	data, err := os.ReadFile(file) // 为文件的数据准备一个变量
	if err == nil {                // 错误处理
		// 如果文件存在
		if err := json.Unmarshal(data, &list); err != nil { // 解析现有数据，整理成一个切片
			fmt.Println()
		}
		// fmt.Println("已读取。")
		return list // 返回这个切片

	} else {
		fmt.Println("读取文件时出错")
		return []AuraClock{} // 出错时返回一个空切片
	}
}

func WriteConfig(list []AuraClock, file string) []AuraClock {

	_, data := fmtJson(list)

	os.WriteFile(file, data, 0o644) // 写入文件
	return list
}

func getSun(lat, lon float64) (time.Time, time.Time) { // 获取日出、日落时间
	// 设定今天的日期
	now := time.Now()
	rise, set := sunrise.SunriseSunset(lat, lon, now.Year(), now.Month(), now.Day()) // 自动计算日出日落时间
	return rise, set
}

func ComClock(aura AuraClock) time.Time { // 通过时间差值和日出、日落时间计算实时计算闹钟时间

	switch {
	case aura.Sun == "sunRise":
		rise := sunRise.Add(aura.Diff) // 加上时间差值
		return rise                    // 格式化闹钟时间

	case aura.Sun == "sunSet":
		set := sunSet.Add(aura.Diff) // 加上时间差值
		return set

	default:
		return time.Now()

	}
}

func FormatClock(aura AuraClock) string {
	tz := latlong.LookupZoneName(latitude, longitude) // 根据经纬度计算时区
	loc, err := time.LoadLocation(tz)                 // 转化为time包的时区
	if err != nil {
		fmt.Println("解析失败，使用默认时区")
		loc = time.Local // 使用默认时区
	}
	return ComClock(aura).In(loc).Format("15:04")
}

func main() {

	// 设置选项
	flag.Float64Var(&longitude, "l", 116.391, "设置经度（默认北京）")
	flag.Float64Var(&latitude, "a", 39.906, "设定纬度（默认北京）")

	flag.StringVar(&ReadPath, "r", "", "指定读取的文件中的闹钟")
	flag.StringVar(&WritePath, "w", "", "创建闹钟并写入指定文件")

	flag.BoolVar(&onlyCreate, "n", false, "只使用创建闹钟的基础功能")

	// flag.BoolVar(&remind, "m", false, "不使用GUI但启动提醒")

	flag.Parse() // 解析参数

	sunRise, sunSet = getSun(latitude, longitude) // 获取日出和日落时间

	w := tabwriter.NewWriter(os.Stderr, 0, 8, 2, ' ', 0) // 创建一个输出接受器
	switch {
	case onlyCreate && ReadPath == "" && WritePath == "": // 只有n选项
		aura := CreateClock()
		fmt.Fprintf(w, "闹钟名称: %s\t| 执行时间: %s\n", aura.Name, FormatClock(aura))

	case !onlyCreate && ReadPath != "" && WritePath == "": // 只有r选项
		auraList = ReadConfig(ReadPath)
		for _, aura := range auraList {
			fmt.Fprintf(w, "闹钟名称: %s\t| 执行时间: %s\n", aura.Name, FormatClock(aura))
		}

	case !onlyCreate && ReadPath == "" && WritePath != "": // 只有w选项
		aura := CreateClock()
		fmt.Fprintf(w, "闹钟名称: %s\t| 执行时间: %s\n", aura.Name, FormatClock(aura))
		auraList = append(auraList, aura)
		ww := WriteConfig(auraList, WritePath)
		fmt.Println("闹钟已写入：", WritePath)
		fmt.Println("文件内容：", ww)

	// case remind:

	default:
		fmt.Println("默认以GUI模式启动；如果使用命令行则选项n、r、w有且只能有一个。")
		StartGUI()
		SortAuraList(&auraList)
		UpdateFile(auraList, "./GUI.json")
		// TriggerClock(auraList[len(auraList)-1])
	}
	w.Flush() // 输出

}
