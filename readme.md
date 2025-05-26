# ⚡ SwiftDork: Your Go-To Dork Query Generator

**English:**
**SwiftDork** is a powerful and flexible Dork query generator written in **Go**. Quickly craft unique Dork queries to discover web vulnerabilities or fuel your research, all at lightning speed!

**Русский:**
**SwiftDork** — это мощный и гибкий генератор Dork-запросов, написанный на **Go**. Быстро создавайте уникальные Dork-запросы для поиска веб-уязвимостей или исследовательских целей, всё это с молниеносной скоростью!

---

## ✨ Key Features / Ключевые возможности

* **Diverse Dictionaries / Разнообразные словари:** Generate queries using admin panels, vulnerabilities, file types, years, locations, and more. / Генерируйте запросы, используя панели администрирования, уязвимости, типы файлов, года, местоположения и многое другое.
* **Targeted Filtering / Целевая фильтрация:** Filter by country and domain for precise results. / Фильтруйте по странам и доменам для получения точных результатов.
* **Multilingual Support / Многоязычная поддержка:** Leverage English and Russian dictionaries, plus a list of popular cities. / Используйте английские и русские словари, а также список популярных городов.
* **Smart Expansion / Умное расширение:** Dictionaries automatically expand with variations for broader searches. / Словари автоматически расширяются за счёт генерации вариаций для более широкого поиска.
* **Custom Templates / Пользовательские шаблоны:** Load your own Dork query templates from `base_templates.txt`. / Загружайте собственные шаблоны Dork-запросов из `base_templates.txt`.
* **Dynamic Data / Динамические данные:** Automatically fetches up-to-date domain lists from IANA and countries from restcountries.com. / Автоматически загружает актуальные списки доменов с IANA и стран с restcountries.com.
* **Blazing Fast / Молниеносная скорость:** Multithreaded generation ensures high performance. / Многопоточная генерация обеспечивает высокую производительность.
* **Specialized Targets / Специализированные цели:** Support for WordPress, Joomla, Nginx, Laravel, and other specific platforms. / Поддержка WordPress, Joomla, Nginx, Laravel и других специфических платформ.
* **Exportable Results / Экспортируемые результаты:** Save your generated queries to a specified file. / Сохраняйте сгенерированные запросы в указанный файл.

---

## 🚀 Get Started / Начало работы

```bash
# Clone the repository / Склонировать репозиторий
git clone https://github.com/helene92727630/swiftdork.git
cd swiftdork

# Build the executable / Скомпилировать исполняемый файл
go build -o swiftdork main.go
