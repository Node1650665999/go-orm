package session

import (
	"errors"
	"geeorm/clause"
	"reflect"
)

// Insert one or more records in database
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		s.CallMethod(BeforeInsert, value)
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	//recordValues 数据格式类似这种: [[Tom 18] [Sam 25]]
	s.clause.Set(clause.VALUES, recordValues...)

	//构造出最终的 SQL 语句和待绑定的 Value
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

// Find gets all eligible records
func (s *Session) Find(values interface{}) error {
	s.CallMethod(BeforeQuery, nil)

	//destSlice 数据格式类似这种: []main.User{}
	destSlice := reflect.Indirect(reflect.ValueOf(values))

	//这里的 Elem() 提取了数组中对象的类型
	destType := destSlice.Type().Elem()

	//获取values对应类型的零值, 数据格式类似这种： main.User{Name:"", Age:0}
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		//dest 数据格式类似这种: dest := User{Name: , Age:0}
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			//获取每个字段的地址，类似这种： values = []interface{}{&u.Name, &u.Age}
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		//相当于传入成员 &u.Name, &u.Age 来获取值
		if err := rows.Scan(values...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, dest.Addr().Interface())

		//相当于将单个 User{} 对象添加进前面的 []User{} 中
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// First gets the 1st row
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	//reflect.New 相当于构造出来了 &[]User{}, 经过Elem() 解指针后得到 []User{}
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	//destSlice.Addr() 又再次构造了 &[]User{}, 绕这一圈的目的在于方便后面直接使用解指针的 destSlice来调用调用方法 Len() 和 Index()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

// Limit adds limit condition to clause
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where adds limit condition to clause
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	//values 格式类似这种: ["name=?,age=?,phone=?", "tom", 18, "1111111"]
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy adds order by condition to clause
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// Update 接受 2 种入参，平铺开来的键值对和 map 类型的键值对
// map: map[string]interface{}  平铺的键值对: "Name", "Tom", "Age", 18, ....
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		//将展开的 kv list 参数("Name", "Tom", "Age", 18 )处理成 map 类型
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	//数据格式类似这种： sql:update user set name=?, age=?;  vars:["Tom",18]
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

// Delete records with where clause
func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

// Count records with where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}
