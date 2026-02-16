package newrule

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sdojjy/genshin-value-rule/pkg/eval"
)

type NewRule struct {
}

func New() *NewRule {
	return &NewRule{}
}

// CharacterInfo 存储角色的价格和专武信息
type CharacterInfo struct {
	Name              string
	Prices            [7]float64 // 0命到6命的价格
	SpecializedWeapon string
}

// WeaponInfo 存储武器的价格信息
type WeaponInfo struct {
	Name   string
	Prices [5]float64 // 精1到精5的价格
}

// RequiredChar 定义了溢价组合中对角色的要求
type RequiredChar struct {
	Name     string
	MinConst int // 最小命座要求
	MaxConst int // 最大命座要求
}

// ComboRule 定义了一条溢价组合规则
type ComboRule struct {
	Name          string
	Value         float64
	RequiredChars []RequiredChar
}

// ValuationRules 包含所有估值规则
type ValuationRules struct {
	Characters map[string]CharacterInfo
	Weapons    map[string]WeaponInfo
	Combos     []ComboRule

	// 角色数量溢价规则
	CharCountMultiplierTiers []struct {
		MinCount int
		MaxCount int
		Factor   float64
	}

	// 资源价值规则
	ResourceValueTiers []struct {
		MinFates int
		Price    float64
	}

	// 特殊规则相关角色列表
	HotC6Chars       []string
	SpecialC2C5Chars []string
}

// 全局变量，存储加载后的所有规则
var rules = loadValuationRules()

// loadValuationRules 初始化所有估值规则，数据源于Word文档
func loadValuationRules() ValuationRules {
	r := ValuationRules{}
	// 角色价格表 [cite: 22, 228]
	r.Characters = map[string]CharacterInfo{
		"杜林":    {Name: "杜林", Prices: [7]float64{5, 10, 50, 60, 70, 200, 600}, SpecializedWeapon: "黑蚀"},
		"伊涅芙":   {Name: "伊涅芙", Prices: [7]float64{5, 10, 30, 35, 40, 100, 500}, SpecializedWeapon: "支离轮光"}, // 600->500
		"丝柯克":   {Name: "丝柯克", Prices: [7]float64{5, 10, 80, 90, 100, 200, 500}, SpecializedWeapon: "苍耀"},  // 550->500
		"爱可菲":   {Name: "爱可菲", Prices: [7]float64{5, 10, 50, 55, 60, 100, 400}, SpecializedWeapon: "香韵奏者"},
		"瓦雷莎":   {Name: "瓦雷莎", Prices: [7]float64{5, 10, 80, 90, 100, 200, 600}, SpecializedWeapon: "溢彩心念"}, // 800->600
		"茜特菈莉":  {Name: "茜特菈莉", Prices: [7]float64{5, 10, 50, 55, 60, 100, 400}, SpecializedWeapon: "祭星者之望"},
		"玛薇卡":   {Name: "玛薇卡", Prices: [7]float64{5, 10, 80, 90, 100, 200, 500}, SpecializedWeapon: "焚曜千阳"}, // 550->500
		"恰斯卡":   {Name: "恰斯卡", Prices: [7]float64{5, 10, 50, 60, 70, 100, 600}, SpecializedWeapon: "星鹫赤羽"},  // 650->600
		"希诺宁":   {Name: "希诺宁", Prices: [7]float64{5, 10, 50, 55, 60, 100, 300}, SpecializedWeapon: "岩峰巡歌"},
		"基尼奇":   {Name: "基尼奇", Prices: [7]float64{5, 10, 50, 60, 70, 100, 380}, SpecializedWeapon: "山王长牙"},
		"玛拉妮":   {Name: "玛拉妮", Prices: [7]float64{5, 10, 50, 60, 70, 100, 380}, SpecializedWeapon: "冲浪时光"},
		"艾梅莉埃":  {Name: "艾梅莉埃", Prices: [7]float64{5, 10, 20, 25, 30, 50, 200}, SpecializedWeapon: "柔灯挽歌"},
		"克洛琳德":  {Name: "克洛琳德", Prices: [7]float64{5, 10, 25, 30, 35, 50, 360}, SpecializedWeapon: "赦罪"},
		"阿蕾奇诺":  {Name: "阿蕾奇诺", Prices: [7]float64{5, 10, 25, 30, 35, 100, 380}, SpecializedWeapon: "赤月之形"},
		"希格雯":   {Name: "希格雯", Prices: [7]float64{5, 10, 15, 20, 25, 50, 200}, SpecializedWeapon: "白雨心弦"},
		"千织":    {Name: "千织", Prices: [7]float64{5, 10, 15, 20, 25, 30, 360}, SpecializedWeapon: "有乐御簾切"},
		"闲云":    {Name: "闲云", Prices: [7]float64{5, 10, 25, 30, 35, 40, 200}, SpecializedWeapon: "鹤鸣余音"},
		"娜维娅":   {Name: "娜维娅", Prices: [7]float64{5, 10, 15, 20, 25, 30, 200}, SpecializedWeapon: "裁断"},
		"芙宁娜":   {Name: "芙宁娜", Prices: [7]float64{5, 10, 30, 35, 40, 80, 250}, SpecializedWeapon: "静水流涌之辉"},
		"那维莱特":  {Name: "那维莱特", Prices: [7]float64{5, 10, 15, 20, 25, 80, 300}, SpecializedWeapon: "万世流涌大典"},
		"莱欧斯利":  {Name: "莱欧斯利", Prices: [7]float64{5, 10, 15, 20, 25, 30, 250}, SpecializedWeapon: "金流监督"},
		"林尼":    {Name: "林尼", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "最初的大魔术"},
		"白术":    {Name: "白术", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "碧落之珑"},
		"艾尔海森":  {Name: "艾尔海森", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "裁叶萃光"},
		"流浪者":   {Name: "流浪者", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "图莱杜拉的回忆"},
		"纳西妲":   {Name: "纳西妲", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "千夜浮梦"},
		"赛诺":    {Name: "赛诺", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "赤沙之杖"},
		"妮露":    {Name: "妮露", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "圣显之钥"},
		"神里绫人":  {Name: "神里绫人", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "波乱月白经津"},
		"申鹤":    {Name: "申鹤", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "息灾"},
		"夜兰":    {Name: "夜兰", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "若水"},
		"八重神子":  {Name: "八重神子", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "神乐之真意"},
		"荒泷一斗":  {Name: "荒泷一斗", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "赤角石溃杵"},
		"珊瑚宫心海": {Name: "珊瑚宫心海", Prices: [7]float64{5, 10, 15, 20, 25, 30, 100}, SpecializedWeapon: "不灭月华"},
		"雷电将军":  {Name: "雷电将军", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "薙草之稻光"},
		"优菈":    {Name: "优菈", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "松籁响起之时"},
		"宵宫":    {Name: "宵宫", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "飞雷之弦振"},
		"枫原万叶":  {Name: "枫原万叶", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "苍古自由之誓"},
		"胡桃":    {Name: "胡桃", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "护摩之杖"},
		"甘雨":    {Name: "甘雨", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "阿莫斯之弓"},
		"达达利亚":  {Name: "达达利亚", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "冬极白星"},
		"钟离":    {Name: "钟离", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "贯虹之槊"},
		"魈":     {Name: "魈", Prices: [7]float64{5, 10, 15, 20, 25, 30, 150}, SpecializedWeapon: "和璞鸢"},
		"可莉":    {Name: "可莉", Prices: [7]float64{5, 10, 15, 20, 25, 30, 100}, SpecializedWeapon: "四风原典"},
		"温迪":    {Name: "温迪", Prices: [7]float64{5, 10, 15, 20, 25, 30, 180}, SpecializedWeapon: "终末嗟叹之诗"},
		"菈乌玛":   {Name: "菈乌玛", Prices: [7]float64{5, 10, 50, 60, 70, 200, 400}, SpecializedWeapon: "纺夜天镜"},
		"菲林斯":   {Name: "菲林斯", Prices: [7]float64{5, 10, 50, 60, 70, 200, 500}, SpecializedWeapon: "血染荒城"},
		"奈芙尔":   {Name: "奈芙尔", Prices: [7]float64{5, 10, 50, 60, 70, 200, 700}, SpecializedWeapon: "真语秘匣"},  // 600->700
		"哥伦比娅":  {Name: "哥伦比娅", Prices: [7]float64{5, 10, 50, 60, 70, 200, 700}, SpecializedWeapon: "帷间夜曲"}, // 新增
		"兹白":    {Name: "兹白", Prices: [7]float64{5, 10, 50, 60, 70, 200, 900}, SpecializedWeapon: "朏魄含光"},   // 新增, C5纠正为200
	}

	// 武器价格表
	r.Weapons = map[string]WeaponInfo{
		"黑蚀":      {Name: "黑蚀", Prices: [5]float64{5, 10, 15, 20, 200}}, // 新增武器
		"支离轮光":    {Name: "支离轮光", Prices: [5]float64{5, 10, 15, 20, 150}},
		"苍耀":      {Name: "苍耀", Prices: [5]float64{5, 10, 15, 20, 200}}, // 300->200
		"香韵奏者":    {Name: "香韵奏者", Prices: [5]float64{5, 10, 15, 20, 150}},
		"溢彩心念":    {Name: "溢彩心念", Prices: [5]float64{5, 10, 15, 20, 250}},
		"祭星者之望":   {Name: "祭星者之望", Prices: [5]float64{5, 10, 15, 20, 150}},
		"焚曜千阳":    {Name: "焚曜千阳", Prices: [5]float64{5, 10, 15, 20, 250}}, // 300->250
		"星鹫赤羽":    {Name: "星鹫赤羽", Prices: [5]float64{5, 10, 15, 20, 250}}, // 300->250
		"岩峰巡歌":    {Name: "岩峰巡歌", Prices: [5]float64{5, 10, 15, 20, 100}},
		"山王长牙":    {Name: "山王长牙", Prices: [5]float64{5, 10, 15, 20, 150}},
		"冲浪时光":    {Name: "冲浪时光", Prices: [5]float64{5, 10, 15, 20, 150}},
		"柔灯挽歌":    {Name: "柔灯挽歌", Prices: [5]float64{5, 10, 15, 20, 50}},
		"赦罪":      {Name: "赦罪", Prices: [5]float64{5, 10, 15, 20, 150}},
		"赤月之形":    {Name: "赤月之形", Prices: [5]float64{5, 10, 15, 20, 200}},
		"白雨心弦":    {Name: "白雨心弦", Prices: [5]float64{5, 10, 15, 20, 50}},
		"有乐御簾切":   {Name: "有乐御簾切", Prices: [5]float64{5, 10, 15, 20, 150}},
		"鹤鸣余音":    {Name: "鹤鸣余音", Prices: [5]float64{5, 10, 15, 20, 50}},
		"裁断":      {Name: "裁断", Prices: [5]float64{5, 10, 15, 20, 50}},
		"静水流涌之辉":  {Name: "静水流涌之辉", Prices: [5]float64{5, 10, 15, 20, 100}},
		"万世流涌大典":  {Name: "万世流涌大典", Prices: [5]float64{5, 10, 15, 20, 150}},
		"金流监督":    {Name: "金流监督", Prices: [5]float64{5, 10, 15, 20, 80}},
		"最初的大魔术":  {Name: "最初的大魔术", Prices: [5]float64{5, 10, 15, 20, 50}},
		"碧落之珑":    {Name: "碧落之珑", Prices: [5]float64{5, 10, 15, 20, 50}},
		"裁叶萃光":    {Name: "裁叶萃光", Prices: [5]float64{5, 10, 15, 20, 25}},
		"图莱杜拉的回忆": {Name: "图莱杜拉的回忆", Prices: [5]float64{5, 10, 15, 20, 25}},
		"千夜浮梦":    {Name: "千夜浮梦", Prices: [5]float64{5, 10, 15, 20, 25}},
		"赤沙之杖":    {Name: "赤沙之杖", Prices: [5]float64{5, 10, 15, 20, 25}},
		"圣显之钥":    {Name: "圣显之钥", Prices: [5]float64{5, 10, 15, 20, 25}},
		"波乱月白经津":  {Name: "波乱月白经津", Prices: [5]float64{5, 10, 15, 20, 25}},
		"息灾":      {Name: "息灾", Prices: [5]float64{5, 10, 15, 20, 25}},
		"若水":      {Name: "若水", Prices: [5]float64{5, 10, 15, 20, 25}},
		"神乐之真意":   {Name: "神乐之真意", Prices: [5]float64{5, 10, 15, 20, 25}},
		"赤角石溃杵":   {Name: "赤角石溃杵", Prices: [5]float64{5, 10, 15, 20, 25}},
		"不灭月华":    {Name: "不灭月华", Prices: [5]float64{5, 10, 15, 20, 25}},
		"薙草之稻光":   {Name: "薙草之稻光", Prices: [5]float64{5, 10, 15, 20, 25}},
		"松籁响起之时":  {Name: "松籁响起之时", Prices: [5]float64{5, 10, 15, 20, 25}},
		"飞雷之弦振":   {Name: "飞雷之弦振", Prices: [5]float64{5, 10, 15, 20, 25}},
		"苍古自由之誓":  {Name: "苍古自由之誓", Prices: [5]float64{5, 10, 15, 20, 25}},
		"护摩之杖":    {Name: "护摩之杖", Prices: [5]float64{5, 10, 15, 20, 25}},
		"阿莫斯之弓":   {Name: "阿莫斯之弓", Prices: [5]float64{5, 5, 5, 5, 25}},
		"冬极白星":    {Name: "冬极白星", Prices: [5]float64{5, 10, 15, 20, 25}},
		"贯虹之槊":    {Name: "贯虹之槊", Prices: [5]float64{5, 5, 5, 5, 25}},
		"和璞鸢":     {Name: "和璞鸢", Prices: [5]float64{5, 5, 5, 5, 25}},
		"四风原典":    {Name: "四风原典", Prices: [5]float64{5, 5, 5, 5, 25}},
		"终末嗟叹之诗":  {Name: "终末嗟叹之诗", Prices: [5]float64{5, 5, 5, 5, 25}},
		"纺夜天镜":    {Name: "纺夜天镜", Prices: [5]float64{5, 10, 15, 20, 100}},
		"血染荒城":    {Name: "血染荒城", Prices: [5]float64{5, 10, 15, 20, 250}},
		"真语秘匣":    {Name: "真语秘匣", Prices: [5]float64{5, 10, 15, 20, 200}},
		"帷间夜曲":    {Name: "帷间夜曲", Prices: [5]float64{5, 10, 15, 20, 200}}, // 新增
		"朏魄含光":    {Name: "朏魄含光", Prices: [5]float64{5, 10, 15, 20, 250}}, // 新增
	}

	// 完整溢价组合 [cite: 25-175, 177-184]
	// 完整溢价组合 [cite: 25-175, 177-184]
	r.Combos = getCombos()

	sort.Slice(r.Combos, func(i, j int) bool {
		return r.Combos[i].Value > r.Combos[j].Value
	})

	// 角色数量乘数规则（未变动）
	r.CharCountMultiplierTiers = []struct {
		MinCount int
		MaxCount int
		Factor   float64
	}{
		{0, 10, 0.6}, {11, 20, 0.8}, {21, 39, 1.0},
		{40, 45, 1.2}, {46, 50, 1.4}, {51, 999, 1.6},
	}

	// 资源价值规则 [cite: 396-404]
	r.ResourceValueTiers = []struct {
		MinFates int
		Price    float64
	}{
		{1000, 1.8},
		{900, 1.7},
		{800, 1.6},
		{700, 1.4},
		{600, 1.3},
		{500, 1.2},
		{300, 1.0},
		{200, 0.5},
	}

	// 特殊规则相关角色列表
	// 特殊规则相关角色列表
	r.HotC6Chars = []string{"杜林", "奈芙尔", "菈乌玛", "菲林斯", "基尼奇", "千织", "瓦雷莎", "克洛琳德", "玛拉妮", "哥伦比娅", "兹白"}

	return r
}

// CalculateValuation 是估值的主入口函数
func (n *NewRule) CalculateValuation(account eval.Assets) eval.ValuationResult {
	var sb strings.Builder

	// --- 步骤一: 计算最优溢价组合 ---
	fmt.Fprintf(&sb, "<div class='step'><h3>步骤一: 计算最优溢价组合附加价值</h3><pre>")
	satisfiedCombos := findSatisfiedCombos(account)
	bestComboBonus, bestComboSelection := findBestComboSelection(satisfiedCombos, make(map[string]bool))
	if len(bestComboSelection) > 0 {
		fmt.Fprintf(&sb, "命中以下最优组合方案，获得附加价值: %.2f\n", bestComboBonus)
		for _, combo := range bestComboSelection {
			fmt.Fprintf(&sb, "  - %s (附加 %.2f)\n", combo.Name, combo.Value)
		}
	} else {
		sb.WriteString("未命中任何溢价组合。\n")
	}
	sb.WriteString("</pre></div>")

	// --- 步骤二: 区分并计算基础价值 ---
	fmt.Fprintf(&sb, "<div class='step'><h3>步骤二: 计算并区分角色与武器的基础价值</h3><pre>")
	applicableValue, exemptValue, baseValueBreakdown := calculateBaseValue(account, bestComboSelection)
	sb.WriteString(baseValueBreakdown)
	fmt.Fprintf(&sb, "\n&gt;&gt; 适用乘数的基础价值: %.2f\n", applicableValue)
	fmt.Fprintf(&sb, "&gt;&gt; 豁免乘数的基础价值: %.2f\n", exemptValue)
	sb.WriteString("</pre></div>")

	// --- 步骤三: 应用角色数量乘数 ---
	fmt.Fprintf(&sb, "<div class='step'><h3>步骤三: 对适用部分应用角色数量乘数</h3><pre>")
	adjustedApplicableValue, multiplierBreakdown := applyCharacterCountMultiplier(applicableValue, len(account.Characters))
	sb.WriteString(multiplierBreakdown)
	sb.WriteString("</pre></div>")

	// --- 步骤四: 计算总基础价值 ---
	totalAdjustedBaseValue := adjustedApplicableValue + exemptValue
	fmt.Fprintf(&sb, "<div class='subtotal'><p>总基础价值 = 调整后适用价值 (%.2f) + 豁免价值 (%.2f) = <strong>%.2f</strong></p></div>", adjustedApplicableValue, exemptValue, totalAdjustedBaseValue)

	// --- 步骤五: 计算资源价值 ---
	fmt.Fprintf(&sb, "<div class='step'><h3>步骤五: 计算资源价值</h3><pre>")
	resourceValue, resourceBreakdown := calculateResourceValue(account)
	sb.WriteString(resourceBreakdown)
	sb.WriteString("</pre></div>")

	// --- 步骤六: 应用特殊规则增益 ---
	fmt.Fprintf(&sb, "<div class='step'><h3>步骤六: 应用特殊规则附加增益</h3><pre>")
	specialBonus, specialBonusBreakdown := applySpecialRules(account, bestComboSelection)
	sb.WriteString(specialBonusBreakdown)
	sb.WriteString("</pre></div>")

	// --- 步骤七: 最终合计 ---
	totalValue := totalAdjustedBaseValue + bestComboBonus + resourceValue + specialBonus
	fmt.Fprintf(&sb, "<div class='final-total'><h3>最终合计</h3>")
	fmt.Fprintf(&sb, "<p>总基础价值    : %.2f</p>", totalAdjustedBaseValue)
	fmt.Fprintf(&sb, "<p>组合附加价值  : %.2f</p>", bestComboBonus)
	fmt.Fprintf(&sb, "<p>资源价值      : %.2f</p>", resourceValue)
	fmt.Fprintf(&sb, "<p>特殊规则增益  : %.2f</p>", specialBonus)
	fmt.Fprintf(&sb, "<hr><p><strong>账号总估值: %.2f</strong></p>", totalValue)
	fmt.Fprintf(&sb, "</div>")

	return eval.ValuationResult{
		FinalTotal: totalValue,
		Breakdown:  sb.String(),
	}
}

// calculateBaseValue 区分计算适用和豁免乘数的基础价值
func calculateBaseValue(account eval.Assets, bestRules []ComboRule) (applicableValue float64, exemptValue float64, breakdown string) {
	var sb strings.Builder

	premiumC6Chars := make(map[string]bool)
	for _, combo := range bestRules {
		for _, req := range combo.RequiredChars {
			if req.MinConst == 6 {
				if c, ok := account.Characters[req.Name]; ok && c == 6 {
					premiumC6Chars[req.Name] = true
				}
			}
		}
	}
	if len(bestRules) > 0 {
		for _, hotChar := range rules.HotC6Chars {
			if c, ok := account.Characters[hotChar]; ok && c == 6 {
				premiumC6Chars[hotChar] = true
			}
		}
	}

	c6CharWeapons := make(map[string]bool)
	for name, constellation := range account.Characters {
		if constellation == 6 {
			if charInfo, ok := rules.Characters[name]; ok && charInfo.SpecializedWeapon != "" {
				c6CharWeapons[charInfo.SpecializedWeapon] = true
			}
		}
	}

	charNames := make([]string, 0, len(account.Characters))
	for name := range account.Characters {
		charNames = append(charNames, name)
	}
	sort.Strings(charNames)

	for _, name := range charNames {
		constellation := account.Characters[name]
		charInfo, ok := rules.Characters[name]
		if !ok {
			continue
		}
		value := charInfo.Prices[constellation]
		reason := ""
		if constellation >= 2 && constellation <= 6 {
			if _, hasWeapon := account.Weapons[charInfo.SpecializedWeapon]; !hasWeapon {
				value *= 0.8
				reason = " (无专武, 8折)"
			}
		}

		if constellation == 6 {
			exemptValue += value
			sb.WriteString(fmt.Sprintf("  - [豁免] 角色 [%s %d命]: %.2f%s\n", name, constellation, value, reason))
		} else {
			applicableValue += value
			sb.WriteString(fmt.Sprintf("  - [适用] 角色 [%s %d命]: %.2f%s\n", name, constellation, value, reason))
		}
	}

	weaponNames := make([]string, 0, len(account.Weapons))
	for name := range account.Weapons {
		weaponNames = append(weaponNames, name)
	}
	sort.Strings(weaponNames)

	for _, name := range weaponNames {
		refine := account.Weapons[name]
		weaponInfo, ok := rules.Weapons[name]
		if !ok {
			continue
		}
		if refine <= 0 {
			refine = 1
		}
		value := weaponInfo.Prices[refine-1]
		reason := ""

		ownerName := ""
		for charName, charInfo := range rules.Characters {
			if charInfo.SpecializedWeapon == name {
				ownerName = charName
				break
			}
		}

		if refine == 5 {
			if ownerName != "" {
				if ownerConst, hasOwner := account.Characters[ownerName]; !hasOwner || ownerConst < 6 {
					value = weaponInfo.Prices[3] // 按精4计价
					reason = fmt.Sprintf(" (角色%s非6命, 按精4计价)", ownerName)
				}
			}
		}

		if refine == 5 && ownerName != "" && premiumC6Chars[ownerName] {
			value *= 2
			reason += " (命中组合内6命角色专武, 价格x2)"
		}

		if c6CharWeapons[name] {
			exemptValue += value
			sb.WriteString(fmt.Sprintf("  - [豁免] 武器 [%s 精%d]: %.2f%s\n", name, refine, value, reason))
		} else {
			applicableValue += value
			sb.WriteString(fmt.Sprintf("  - [适用] 武器 [%s 精%d]: %.2f%s\n", name, refine, value, reason))
		}
	}

	if sb.Len() == 0 {
		return 0, 0, "账号内无有效角色或武器。\n"
	}
	return applicableValue, exemptValue, sb.String()
}

// findSatisfiedCombos 找出账号满足的所有组合
func findSatisfiedCombos(account eval.Assets) []ComboRule {
	var satisfied []ComboRule
	for _, combo := range rules.Combos {
		isSatisfied := true
		for _, req := range combo.RequiredChars {
			constellation, ok := account.Characters[req.Name]
			if !ok || constellation < req.MinConst || constellation > req.MaxConst {
				isSatisfied = false
				break
			}
		}
		if isSatisfied {
			satisfied = append(satisfied, combo)
		}
	}
	return satisfied
}

// findBestComboSelection 使用回溯算法寻找最优组合的附加价值
func findBestComboSelection(availableCombos []ComboRule, usedChars map[string]bool) (float64, []ComboRule) {
	if len(availableCombos) == 0 {
		return 0, nil
	}
	// Case 1: Skip the current combo
	valueA, selectionA := findBestComboSelection(availableCombos[1:], usedChars)

	// Case 2: Try to select the current combo
	currentCombo := availableCombos[0]
	canSelect := true
	for _, req := range currentCombo.RequiredChars {
		if usedChars[req.Name] {
			canSelect = false
			break
		}
	}

	if canSelect {
		newUsedChars := make(map[string]bool)
		for k, v := range usedChars {
			newUsedChars[k] = v
		}
		for _, req := range currentCombo.RequiredChars {
			newUsedChars[req.Name] = true
		}
		valueB, selectionB := findBestComboSelection(availableCombos[1:], newUsedChars)
		valueB += currentCombo.Value
		selectionB = append(selectionB, currentCombo)

		if valueB > valueA {
			return valueB, selectionB
		}
	}

	return valueA, selectionA
}

// calculateResourceValue 计算资源价值
func calculateResourceValue(account eval.Assets) (float64, string) {
	totalFates := account.JiuChanZhiYuan + (account.YuanShi / 160)
	var sb strings.Builder
	fmt.Fprintf(&sb, "账号总资源: %d 原石 + %d 纠缠之源 = %d 总抽数\n", account.YuanShi, account.JiuChanZhiYuan, totalFates)

	if totalFates < 200 {
		sb.WriteString("总抽数低于200，不计价。\n")
		return 0, sb.String()
	}

	value := 0.0
	for _, tier := range rules.ResourceValueTiers {
		if totalFates >= tier.MinFates {
			value = float64(totalFates) * tier.Price
			fmt.Fprintf(&sb, "  - %d 抽: %d * %.2f = %.2f\n",
				totalFates, totalFates, tier.Price, value)
			break
		}
	}
	fmt.Fprintf(&sb, "资源总价值: %.2f\n", value)
	return value, sb.String()
}

// applyCharacterCountMultiplier 应用角色数量乘数
func applyCharacterCountMultiplier(applicableValue float64, charCount int) (float64, string) {
	for _, tier := range rules.CharCountMultiplierTiers {
		if charCount >= tier.MinCount && charCount <= tier.MaxCount {
			finalValue := applicableValue * tier.Factor
			return finalValue, fmt.Sprintf("账号有 %d 个五星角色，对适用部分应用 %.0f%% 的乘数:\n  %.2f * %.2f = %.2f\n", charCount, tier.Factor*100, applicableValue, tier.Factor, finalValue)
		}
	}
	return applicableValue, fmt.Sprintf("账号有 %d 个五星角色，未找到对应的乘数规则，价值不变。\n", charCount)
}

// applySpecialRules 应用特殊规则增益
func applySpecialRules(account eval.Assets, bestRules []ComboRule) (float64, string) {
	var totalBonus float64
	var sb strings.Builder

	for _, combo := range bestRules {
		isMaxConstellationCombo := false
		for _, req := range combo.RequiredChars {
			if req.MinConst == 6 {
				isMaxConstellationCombo = true

				break
			}
		}

		if isMaxConstellationCombo {
			// Rule: C2-C5 bonus per combo
			specialCharsFoundInThisCombo := 0
			for _, specialChar := range rules.SpecialC2C5Chars {
				for _, req := range combo.RequiredChars {
					if req.Name == specialChar {
						if c, inAccount := account.Characters[specialChar]; inAccount && c >= 2 && c <= 5 {
							specialCharsFoundInThisCombo++
						}
						break
					}
				}
			}
			if specialCharsFoundInThisCombo >= 3 {
				totalBonus += 400
				fmt.Fprintf(&sb, "  - 组合 [%.30s...] 包含3种特定2-5命角色，附加价值 +400\n", combo.Name)
			} else if specialCharsFoundInThisCombo == 2 {
				totalBonus += 200
				fmt.Fprintf(&sb, "  - 组合 [%.30s...] 包含2种特定2-5命角色，附加价值 +200\n", combo.Name)
			}
		}
	}

	if len(bestRules) > 0 {
		// Rule: Hot C6 Characters [cite: 386-387]
		for _, hotChar := range rules.HotC6Chars {
			if constellation, ok := account.Characters[hotChar]; ok && constellation == 6 {
				totalBonus += 250
				fmt.Fprintf(&sb, "  - 命中热门6命角色 [%s]，附加价值 +250\n", hotChar)

			}
		}
	}

	if sb.Len() == 0 {
		return 0, "未触发任何特殊角色规则。\n"
	}

	return totalBonus, sb.String()
}

func main() {
	// 示例账号，用于演示C6豁免规则和HTML输出
	exampleAccount := eval.Assets{
		Characters: map[string]int{
			"芙宁娜": 6, // 6命角色，其价值应豁免乘数
			"胡桃":  1, // 1命角色，其价值适用乘数
		},
		Weapons: map[string]int{
			"静水流涌之辉": 1, // 芙宁娜专武，其价值应豁免乘数
			"护摩之杖":   1, // 胡桃专武，其价值适用乘数
		},
		YuanShi:        160 * 550, // 550抽
		JiuChanZhiYuan: 0,
	}

	// 执行估值
	rule := New()
	result := rule.CalculateValuation(exampleAccount)

	// 打印估值报告 (HTML格式)
	fmt.Println("============== 账号估值报告 (HTML输出) ==============")
	fmt.Println(result.Breakdown)
	fmt.Printf("\n============== 最终估值: %.2f 元 ==============\n", result.FinalTotal)
}
