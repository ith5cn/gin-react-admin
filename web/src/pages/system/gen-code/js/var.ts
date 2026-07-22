export const realtionsType = [
  { name: '一对一', value: 'hasOne' },
  { name: '一对多', value: 'hasMany' },
  { name: '一对一（反向)', value: 'belongsTo' },
  { name: '多对多', value: 'belongsToMany' }
]

export const queryType = [
  { label: '=', value: 'eq' },
  { label: '!=', value: 'neq' },
  { label: '>', value: 'gt' },
  { label: '>=', value: 'gte' },
  { label: '<', value: 'lt' },
  { label: '<=', value: 'lte' },
  { label: 'LIKE', value: 'like' },
  { label: 'IN', value: 'in' },
  { label: 'NOT IN', value: 'notin' },
  { label: 'BETWEEN', value: 'between' }
]
