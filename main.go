package main

import (
	"fmt"
	"sync"
)

type Character struct {
	Name      string  // название персонажа
	HP        int     // здоровье
	MaxHP     int     // макс. здорове
	Mana      int     // мана
	MaxMana   int     // макс. мана
	Strength  int     // сила
	Agility   int     // ловкость
	Defense   int     // защита
	Weapon    *Weapon // оружие
	Inventory []Item  // инвентарь
	mu        sync.Mutex
}

type Weapon struct {
	Name    string // название оружия
	Damage  int    // урон
	Durable int    // прочность
}

type Item struct {
	Name     string // название итема
	Effect   string // эффект
	Quantity int    // количество
}

type Combatant interface {
	Attack(target *Character) error
	TakeDamage(damage int)
	IsAlive() bool
}

type Usable interface {
	Use(target *Character) error
}

type GameConfig struct {
	StrengthCoefficient float64 // модификатор силы
	DefenseReduction    float64 // модификатор защиты
	CritMultiplier      float64 // модификатор крита
}

var defaultConfig = GameConfig{
	StrengthCoefficient: 0.12,
	DefenseReduction:    0.3,
	CritMultiplier:      2.0,
}

// проверка на жив ли персонаж
func (c *Character) IsAlive() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.HP > 0 {
		fmt.Printf("Здоровье %v: %d\n", c.Name, c.HP)
		return true
	} else {
		fmt.Printf("Персонаж %v убит \n", c.Name)
		return false
	}
}

// атака
func (c *Character) Attack(target *Character) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	// проверка прочности оружия
	if c.Weapon.Durable != 0 {
		// формула расчета урона ((сила * 0.2) * урон)
		Power := (float64(c.Strength) * defaultConfig.StrengthCoefficient) * float64(c.Weapon.Damage)
		c.Weapon.Durable -= 1
		fmt.Printf("Персонаж %v атаковал персонажа %v на %.2f\n", c.Name, target.Name, Power)
		return int(Power)
	}
	fmt.Printf("Оружие сломано!")

	return 0
}

// получение урона
func (c *Character) TakeDamage(damage int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// расчет получнения урона
	// формула получения блокировки урона
	blockDamage := float64(damage) * (float64(c.Defense) * 0.01) // 0.01 = 1% за единицу защиты
	// проверка если заблокированный превышает обычный урон

	if blockDamage > float64(damage) {
		blockDamage = float64(damage)
	}
	fmt.Printf("%v заблокирован урон: %d \n", c.Name, int(blockDamage))

	damage -= int(blockDamage)
	c.HP -= damage
}

// использование предмета
func (c *Character) UseItem(itemName string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, item := range c.Inventory {
		if itemName == item.Name && item.Quantity > 0 { // проверка по имени и по кол-ву
			c.HP += 20

			if c.HP > c.MaxHP {
				c.HP = c.MaxHP
			}
			c.Inventory[i].Quantity--

			fmt.Printf("Персонаж %v успешно вылечился с помощью %v. Теперь его здоровье: %d \n",
				c.Name, item.Name, c.HP)
			return nil
		}
	}
	return fmt.Errorf("Такого предмета в инвентаре нет")
}

func main() {
	hero := &Character{
		Name:     "Джон",
		HP:       100,
		MaxHP:    100,
		Strength: 10,
		Agility:  10,
		Defense:  10,
		Weapon:   &Weapon{"Sword", 30, 50},
		Inventory: []Item{
			{Name: "Зелье лечения", Effect: "Лечение", Quantity: 2},
			{Name: "Яблоко", Effect: "Лечение", Quantity: 2},
		},
	}

	monster := &Character{
		Name:     "Дракон",
		HP:       100,
		Strength: 10,
		Defense:  10,
		Weapon:   &Weapon{"Когти", 20, 50},
	}

	for hero.IsAlive() && monster.IsAlive() {
		fmt.Println("\n1. Атаковать")
		fmt.Println("2. Использовать предмет")
		fmt.Print("Выбери действие: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			damage := hero.Attack(monster)
			monster.TakeDamage(damage)
		case 2:
			fmt.Print("Введи название предмета: ")
			var itemName string
			fmt.Scan(&itemName)
			if err := hero.UseItem(itemName); err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Предмет использован")
			}
		}

		if monster.IsAlive() {
			damage := monster.Attack(hero)
			hero.TakeDamage(damage)
		}

	}
}
