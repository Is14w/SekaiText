export interface CharacterInfo {
  index: number
  name: string
  nameJ: string
  nameC: string
}

export interface UnitInfo {
  key: string
  name: string
}

export const characterDict: CharacterInfo[] = [
  { index: 0, name: 'ichika', nameJ: '一歌', nameC: '一歌' },
  { index: 1, name: 'saki', nameJ: '咲希', nameC: '咲希' },
  { index: 2, name: 'honami', nameJ: '穂波', nameC: '穗波' },
  { index: 3, name: 'shiho', nameJ: '志歩', nameC: '志步' },
  { index: 4, name: 'minori', nameJ: 'みのり', nameC: '实乃里' },
  { index: 5, name: 'haruka', nameJ: '遥', nameC: '遥' },
  { index: 6, name: 'airi', nameJ: '愛莉', nameC: '爱莉' },
  { index: 7, name: 'shizuku', nameJ: '雫', nameC: '雫' },
  { index: 8, name: 'kohane', nameJ: 'こはね', nameC: '心羽' },
  { index: 9, name: 'an', nameJ: '杏', nameC: '杏' },
  { index: 10, name: 'akito', nameJ: '彰人', nameC: '彰人' },
  { index: 11, name: 'touya', nameJ: '冬弥', nameC: '冬弥' },
  { index: 12, name: 'tsukasa', nameJ: '司', nameC: '司' },
  { index: 13, name: 'emu', nameJ: 'えむ', nameC: '笑梦' },
  { index: 14, name: 'nene', nameJ: '寧々', nameC: '宁宁' },
  { index: 15, name: 'rui', nameJ: '類', nameC: '类' },
  { index: 16, name: 'kanade', nameJ: '奏', nameC: '奏' },
  { index: 17, name: 'mafuyu', nameJ: 'まふゆ', nameC: '真冬' },
  { index: 18, name: 'ena', nameJ: '絵名', nameC: '绘名' },
  { index: 19, name: 'mizuki', nameJ: '瑞希', nameC: '瑞希' },
  { index: 20, name: 'miku', nameJ: 'ミク', nameC: 'MIKU' },
  { index: 21, name: 'rin', nameJ: 'リン', nameC: 'RIN' },
  { index: 22, name: 'len', nameJ: 'レン', nameC: 'LEN' },
  { index: 23, name: 'luka', nameJ: 'ルカ', nameC: 'LUKA' },
  { index: 24, name: 'meiko', nameJ: 'MEIKO', nameC: 'MEIKO' },
  { index: 25, name: 'kaito', nameJ: 'KAITO', nameC: 'KAITO' },
  { index: 26, name: 'miku_band', nameJ: 'ミク_LeoN', nameC: 'MIKU' },
  { index: 27, name: 'miku_idol', nameJ: 'ミク_MMJ', nameC: 'MIKU' },
  { index: 28, name: 'miku_street', nameJ: 'ミク_VBS', nameC: 'MIKU' },
  { index: 29, name: 'miku_park', nameJ: 'ミク_WS', nameC: 'MIKU' },
  { index: 30, name: 'miku_nothing', nameJ: 'ミク_25', nameC: 'MIKU' },
]

export function findCharacterByJapaneseName(name: string): CharacterInfo | undefined {
  return characterDict.find(c => c.nameJ === name)
}
