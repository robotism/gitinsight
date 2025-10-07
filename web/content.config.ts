import path from "path";
import { defineContentConfig, defineCollection, z } from "@nuxt/content";

export default defineContentConfig({
  collections: {
    articles: defineCollection({
      type: "page",
      source: {
        cwd: path.resolve("public/static/articles"),
        include: "**/*.md",
      },
      schema: z.object({
        title: z.string(),
        author: z.array(z.string()),
        category: z.array(z.string()),
        tags: z.array(z.string()),
        cover: z.string(),
        createdAt: z.date(),
        updatedAt: z.date(),
        route: z.string(),
        excerpt: z.object({
          type: z.string(),
          children: z.any(),
        }),
      }),
    }),
    projects: defineCollection({
      type: "page",
      source: {
        cwd: path.resolve("public/static/projects"),
        include: "**/*.md",
      },
      schema: z.object({
        title: z.string(),
        author: z.array(z.string()),
        category: z.array(z.string()),
        tags: z.array(z.string()),
        cover: z.string(),
        createdAt: z.date(),
        updatedAt: z.date(),
        excerpt: z.object({
          type: z.string(),
          children: z.any(),
        }),
        route: z.string(),
      }),
    }),
    about: defineCollection({
      type: "page",
      source: {
        cwd: path.resolve("public/static/about"),
        include: "**/*.md",
      },
      schema: z.object({
        title: z.string(),
        author: z.array(z.string()),
        category: z.array(z.string()),
        tags: z.array(z.string()),
        cover: z.string(),
        createdAt: z.date(),
        updatedAt: z.date(),
        excerpt: z.object({
          type: z.string(),
          children: z.any(),
        }),
      }),
    }),
  },
});
