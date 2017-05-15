Что это?
========

Этот репозиторий служит для распространения драйвера для
[Систем безопасности "Астра"](http://www.teko.biz/), для языка программирования
Go.

Драйвер распространяется только в бинарном виде, в связи с тем, что протокол
взаимодействия проприетарный. Последняя версия драйвера всегда доступна для
загрузки в разделе
[Releases](https://github.com/andrey-yantsen/teko-astra-go/releases/latest).

Помимо самого драйвера в архивах распространяются исходный код зависимостей:
[goburrow/serial](https://github.com/goburrow/serial) и
[npat-efault/crc16](https://github.com/npat-efault/crc16)

Перед началом работы
====================

Устройства Астра-РИ-М выходят с завода в режиме "автономный", он не
предполагает работу в качестве ведомого устройства на линии RS-485, в связи с
этим необходимо переключить РИ-М в режим "системный". Делается это заменой
прошивки, для чего потребуется загрузить и установить
[ПО ПКМ Астра Pro](http://www.teko.biz/support/programms/pc/) (во время
установки убедитесь, что среди модулей ПКМ активирована галочка рядом с модулем
смены ПО), и далее следуйте документации из [Инструкции пользователя](http://www.teko.biz/upload/rukovod/RR-RI-M_%D0%98%D0%BD%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%86%D0%B8%D1%8F%20%D0%BF%D0%BE%D0%BB%D1%8C%D0%B7%D0%BE%D0%B2%D0%B0%D1%82%D0%B5%D0%BB%D1%8F.pdf),
раздел "Смена ПО на РР для работы в режиме «Системный»"

Как этим пользоваться?
======================

Go, начиная с версии 1.7, поддерживает `binary-only packages`, это зависимости
распространяемые только в виде бинарных файлов, с минимальными исходниками. При
этом версия Go, используемая вами, и использованная для компиляции пакета,
должны совпадать.

В разделе загрузок можно увидеть множество файлов, например
`go1.8.1_android_386.tar.gz`, — в названии зашифрована версия Go, целевая ОС и
архитектура. Вам необходимо выбрать подходящий вам архив и распаковать его в
`GOPATH`.

В файле [Makefile](Makefile) можно увидеть пример загрузки и установки
корректных версий зависимостей, а в [cmd/astra/main.go](cmd/astra/main.go) —
пример исполняемого файла, взаимодействующего с Астра-РИ-М.

В целом процесс взаимодействия должен строиться следующим образом:
1. Инициализация драйвера командой `driver := astra_l.Connect("/dev/tty1")`, где
`/dev/tty1` — порт, к которому подключен РИ-М
2. Затем необходимо указать адрес устройства, с которым будет проходить
дальнейшая работа: `device := driver.GetDevice(0x01)`
3. `FindDevice()` с адресом устройства `0xFF`, чтобы убедиться что РИ-М
доступен и не зарегистрирован
4. `RegisterDevice()` с аргументом `f.EUI.GetShortDeviceEUI()`, где `f` —
результат вызова `FindDevice()`. После этого шага устройство перестанет отвечать
по адресу `0xFF`, и будет отзываться только на адрес указанный в шаге 2
5. Затем необходимо инициализировать радио-сеть: `device.CreateLevel2Net(1)`,
где аргумент соответствует литере радио-извещателей, с которыми необходимо
работать. Т.е. `1` в качестве аргумента означает что работа будет идти с
радио-извещателями с пометкой `лит.1`
6. Из необходимых радио-извещателей Астра извлечь батарейки, для каждого
регистрируемого извещателя вызвать функцию `device.RegisterLevel2Device(0)`,
разместить извещатель недалеко от РИ-М и установить элемент питания. Если
регистрация будет успешной — функция в ответ вернёт краткую информацию о новом
извещателе.
7. При восстановлении связи с РИ-М желательно запрашивать состояние устройства
вызовом `device.GetState()`
8. Для получения списка событий необходимо вызвать `device.GetEvents()`.
