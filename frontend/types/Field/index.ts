enum F {
  /** 姓名 */ NAME = 'NAME',
  /** 考号 */ ID = 'ID',
  /** 学校 */ SCHOOL = 'SCHOOL',
  /** 班级 */ CLASS = 'CLASS',

  /** 总分 */ SCORED = 'SCORED',
  /** 排名 */ RANK = 'RANK',

  /** 语文 */ YW = 'YW',
  /** 数学 */ SX = 'SX',
  /** 英语 */ YY = 'YY',

  /** 物理 */ WL = 'WL',
  /** 化学 */ HX = 'HX',
  /** 生物 */ SW = 'SW',

  /** 政治 */ ZZ = 'ZZ',
  /** 历史 */ LS = 'LS',
  /** 地理 */ DL = 'DL',

  /** 主科 (语+数+英) */
  ZK = 'ZK',
  /** 理综 (物+化+生) */
  LZ = 'LZ',
  /** 文综 (政+历+地) */
  WZ = 'WZ',

  /** 理科 (语数英+理综) */
  LK = 'LK',
  /** 文科（语数英+理综） */
  WK = 'WK',

  /** 校排名 */
  SCHOOL_RANK = 'SCHOOL_RANK',
  /** 班排名 */
  CLASS_RANK = 'CLASS_RANK',
}

type _ScoreData = {
  [F.NAME]: string, [F.ID]: string,
  [F.SCHOOL]: string,
  [F.CLASS]: string,
  [F.SCORED]: number,
  [F.RANK]: number,
  [F.YW]: number, [F.SX]: number, [F.YY]: number,
  [F.WL]: number, [F.HX]: number, [F.SW]: number,
  [F.ZZ]: number, [F.LS]: number, [F.DL]: number,
  [F.LZ]: number, [F.WZ]: number,
  [F.ZK]: number, [F.LK]: number, [F.WK]: number

  [F.SCHOOL_RANK]: number, [F.CLASS_RANK]: number
}

export default F
export interface ScoreData extends _ScoreData {}
