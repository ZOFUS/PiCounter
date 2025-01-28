# Вычисление числа Пи с использованием Go

Этот проект представляет собой реализацию алгоритма вычисления числа Пи Чудновского с заданной точностью, используя язык программирования Go.

## Описание
 
**Вычисление Пи на CPU:** Использует многопоточную модель для распараллеливания вычислений на CPU.

В проекте также реализовано кэширование факториалов для оптимизации производительности.

## Структура проекта

*   `main.go`: Содержит основной код программы, включая настройку вычислений и вывод результатов.
*   `pi/pi.go`: Содержит логику вычисления числа Пи, включая функции для вычисления факториала, частичных сумм ряда.
*   `pi/pi_test.go`: Содержит модульные тесты для проверки корректности вычислений.

## Зависимости

Для работы проекта необходимы следующие библиотеки:

*   `github.com/schollz/progressbar/v3`: Для отображения прогресса вычислений.
  
Установите пакет, используя `go get`:

```bash
go get github.com/schollz/progressbar/v3
