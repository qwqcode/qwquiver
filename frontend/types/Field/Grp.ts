import F from '.'

/** 主要字段 */
export const F_MAIN = [F.RANK, F.NAME, F.CODE, F.SCHOOL, F.SCHOOL_RANK, F.CLASS, F.CLASS_RANK, F.TOTAL]
export const F_MAIN_SCHOOL = [F.SCHOOL_RANK, F.NAME, F.CODE, F.RANK, F.SCHOOL, F.CLASS, F.CLASS_RANK, F.TOTAL]
export const F_MAIN_CLASS = [F.CLASS_RANK, F.NAME, F.CODE, F.RANK, F.SCHOOL, F.SCHOOL_RANK, F.CLASS, F.TOTAL]

/** 主科 */
export const F_ZK_SUBJ = [F.YW, F.SX, F.YY]
/** 理科 */
export const F_LZ_SUBJ = [F.WL, F.HX, F.SW]
/** 文科 */
export const F_WZ_SUBJ = [F.ZZ, F.LS, F.DL]
/** 科目字段 */
export const F_SUBJ = [...F_ZK_SUBJ, ...F_LZ_SUBJ, ...F_WZ_SUBJ]

/** 拓展求和字段  */
export const F_EXT_SUM = [F.ZK, F.LZ, F.WZ, F.LK, F.WK]

/** 目标排名字段 */
export const F_TARGET_RANK = [F.SCHOOL, F.CLASS, ...F_SUBJ, ...F_EXT_SUM]

/** 数字字段 */
export const F_NUM_ALL = [F.ID, F.RANK, F.SCHOOL_RANK, F.CLASS_RANK, F.TOTAL, ...F_SUBJ, ...F_EXT_SUM]
/** 非数字字段 */
export const F_STR_ALL = [F.NAME, F.CODE, F.SCHOOL, F.CLASS]

export const F_ALL = Object.keys(F) as F[]
