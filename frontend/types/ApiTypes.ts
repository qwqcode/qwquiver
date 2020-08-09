import F from './Field'

/** 数据表配置 */
export interface EXAM_CONF {
  /** 名称 */
  Name: string

  /** 标签 */
  Label?: string

  /** 每科最大分值 */
  FullScore?: { [subject in F]?: number }

  /** 考试日期 */
  Date?: string

  /** 所属分类 */
  Grp?: string

  /** 备注 */
  Note?: string
}

export type EXAM_MAP = EXAM_CONF[]

export interface CommonParms {
  /** 成绩数据表名 */
  exam?: string
}

export interface AllSchoolParams extends CommonParms {}
export interface AllSchoolData {
  school: { [schoolName: string]: string[] }
}

export interface PaginatedData {
  page: number
  pageSize: number
  total: number
  lastPage: number
  list: any[]
}

export interface ConfData {
  examMap: EXAM_MAP
  examGrpList: string[]
  fieldTransDict: {[f in F]?: string}
}

export interface QueryParams extends CommonParms {
  where?: string
  page?: number
  pageSize?: number
  sort?: string
  init?: boolean
}
export interface QueryData extends PaginatedData {
  examName: string
  dataDesc: string
  examConf: EXAM_CONF
  fieldList: F[]
  avgList: { [key in F]?: number }
  condList: { [key in F]?: string|RegExp }
  sortList: { [key in F]?: 1|-1 }
  initConf?: ConfData
}

export interface ChartParams {
  examGrp?: string
  where?: string
}
export interface ChartData {
  uncertain: boolean
  fieldList: string[]
  chartData: {[key in F|'exam']?: any}[]
}
