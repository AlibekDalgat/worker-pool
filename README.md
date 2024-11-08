# Примитивный worker-pool

# Описание задачи
Реализовать на Golang примитивный worker-pool с возможностью динамически добавлять и удалять воркеры. Входные данные (строки) поступают в канал, воркеры их обрабатывают (например, выводят на экран номер воркера и сами данные). Задание на базовые знания каналов и горутин.

# Реализация
Сущности, реализующие работу пула следующие

**Worker** - структура, выполняющая роль воркера. Содержит поля: 
- id - идентификатор воркера,
- workChan - канал, по которому передаются данные (указывает на тот же канал, что и пул),
- killChan - канал, по которому передаётся сигнал о завершении.

Worker имеет функции:
- runWorker(), которая обрабатывает входные данные, переданные в канал пула, и сигналы о завершении канала.

**WorkerPool** - структура, выполняющая роль воркер-пула. Содержит поля:
- workers - множество воркеров
- poolChan - канал, по которому передаются данные
- shutdown - канал, по которому передаётся сигнал о завершении работы воркер-пула
- count - количество работающих воркеров
- wg - счётчик sync.WaitGroup
- mutex - мьютекс sync.Mutex

WorkerPool имеет функции:
- AddWorker - добавиьт воркер
- RemoveWorker - удалить воркер
- Shutdown - завершить работу воркер-пула
- SendTask - отправить данные в канал

Также реализована функция dataProcessing, имитирующая долгую обраюотку данных

# Пример
Пример реализован в функции mai. Для реализации своего примера можно воспользоваться функционалом структур

В ходе примера инициализируется воркер-пул и добавляется три воркера. Далее в отдельной горутине в цикле отправляются числа от 0 до десяти. Пока данные отправляются данные далее происходит спячка на полсекунды (первые три чсила станут обрабатываться всеми имеющимися воркерами), потом удаляется воркер №1, в следствии чего воркер 1 прерывает обработку данных и удаляется, и добавляется ещё один воркер, который станет обрабатывать число 3, и снова спячка. Затем спячка ещё на полсекунды.

За эту секунду обработуются те данные, которые обработались вторым и третьим воркером, этими же воркерами будут приняты для обработки числа 4 и 5. После спячки завершится работа воркер-пула, в следствии чего обработка воркерами  данных 3, 4, 5 прекратиться.
