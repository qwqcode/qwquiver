import F from '.'

const transDictCN = {
  [F.NAME]: '姓名',
  [F.ID]: '考号',
  [F.SCHOOL]: '学校',
  [F.CLASS]: '班级',
  [F.TOTAL]: '总分',
  [F.RANK]: '总排名',
  [F.YW]: '语文', [F.SX]: '数学', [F.YY]: '英语',
  [F.WL]: '物理', [F.HX]: '化学', [F.SW]: '生物',
  [F.ZZ]: '政治', [F.LS]: '历史', [F.DL]: '地理',
  [F.ZK]: '主科', [F.LZ]: '理综', [F.WZ]: '文综',
  [F.LK]: '理科总分', [F.WK]: '文科总分',
  [F.SCHOOL_RANK]: '校排名', [F.CLASS_RANK]: '班排名'
}

export { transDictCN as transDict }

export function FTrans (fieldName: F|string) {
  return transDictCN[fieldName] || ''
}
