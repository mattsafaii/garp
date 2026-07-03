---
title: Blog
layout: base.html
---

{% for post in collections.blog %}* [{{ post.title }}]({{ post.url }}){% if post.date %} — {{ post.date | date:"January 2, 2006" }}{% endif %}
{% endfor %}
