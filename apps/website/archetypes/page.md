---
title: "{{ replace .Name "_" " " | title }}"
description: "" 
date: {{ .Date }}
# url: "/{{ replace .Name "_" " " | title }}"
hidden: false
draft: false
menu:
  main:
    identifier: "{{ replace .Name "_" " " | title }}"
    name: "{{ replace .Name "_" " " | title }}"
    url: "/{{ replace .Name "_" " " | title }}"
    weight: 10
---