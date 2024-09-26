# Design Decisions for `go-htmx`: Choosing Standard Templates over `gomponents` and `templ`

When developing the go-htmx package, it's essential to choose an approach that aligns with the project's goals, leverages Go's strengths, and meets the needs of its users. Here's why we prefer not to adopt the direction of frameworks like gomponents or gotempl:
1. Simplicity and Familiarity
   - **Standard Library Usage**: `go-htmx` utilizes Go's built-in html/template package, which is familiar to most Go developers. This reduces the learning curve and leverages well-established practices.
   - **Avoiding New Abstractions**: Introducing new templating languages or paradigms (as in `gomponents` or `templ`) adds complexity. By sticking with standard templates, we keep things simple and straightforward.

2. Separation of Concerns

   - **Clear Division**: Using templates allows for a clear separation between business logic (Go code) and presentation logic (HTML templates). This promotes cleaner, more maintainable code.
   - **Design Collaboration**: Designers can work on HTML templates without needing to understand Go code, facilitating better collaboration between developers and designers.

3. Readability and Maintainability

   - **Cleaner Codebase**: Mixing HTML generation with Go code can lead to verbose and less readable code. Keeping templates in separate files improves readability.
   - **Ease of Maintenance**: Separate templates are easier to update and maintain, especially in large projects with multiple contributors.

4. Performance Considerations

   - **Efficient Rendering**: Go's template engine is optimized for performance. By using cached templates and avoiding runtime code generation, we achieve efficient rendering.
   - **Reduced Overhead**: Avoiding additional abstraction layers minimizes overhead and potential performance bottlenecks.

5. Flexibility and Extensibility

   - **Template Extensibility**: Standard templates can be extended and customized using template.FuncMap, allowing for powerful template functions without locking into a specific framework's way of doing things.
   - **Partial Rendering Support**: go-htmx supports partial rendering and template wrapping out of the box, providing flexibility in how components are composed.

6. Compatibility and Integration

   - **Ecosystem Compatibility**: By adhering to Go's standard library, go-htmx is compatible with a wide range of existing libraries and tools.
   - **Ease of Integration**: Projects that already use html/template can integrate go-htmx without significant refactoring or adaptation.

7. Avoiding Vendor Lock-in

   - **Long-Term Stability**: Relying on the standard library reduces dependency on third-party packages that may become unmaintained or introduce breaking changes.
   - **Open Standards**: Using widely adopted standards ensures better support and a larger community for troubleshooting and enhancements.

8. Designer-Friendly Approach

   - **Template Editing Tools**: Designers can use standard HTML editors and tools with syntax highlighting, validation, and auto-completion to work on templates.
   - **No Need for Go Expertise**: Keeping templates in HTML allows designers to contribute without needing to learn Go-specific component syntax or code structures.

9. Learning from Other Ecosystems

   - **Avoiding Complexity**: Other languages and frameworks have moved away from embedding logic directly into templates or views due to complexity and maintenance challenges.
   - **Established Best Practices**: The Go community often emphasizes simplicity and clarity. By following established best practices, we align with the community's values.

10. Focused Scope and Purpose

    - **Specific Use Cases**: go-htmx is designed to enhance the rendering of templates in the context of htmx interactions. Adopting a different templating paradigm might dilute its focus.
    - **Maintainability**: Keeping the package focused on its core functionality makes it easier to maintain and evolve over time.

--- 

## Conclusion

While frameworks like `gomponents` and `templ` offer alternative approaches to building web applications in Go, they introduce paradigms that may not align with the goals of the go-htmx package. By leveraging Go's standard templating system and focusing on simplicity, maintainability, and compatibility, go-htmx provides a solution that is both powerful and accessible.

## Summary of Reasons:

- **Simplicity**: Avoiding unnecessary complexity by using standard templates.
- **Maintainability**: Easier to read, update, and maintain codebases.
- **Compatibility**: Seamless integration with existing Go tools and libraries.
- **Performance**: Efficient rendering without additional overhead.
- **Community Alignment**: Following Go community best practices and conventions.

Note: It's important to choose the right tool for the job. While alternative frameworks may be suitable for certain projects, the decision to stick with standard templates in go-htmx is based on the desire to keep the package simple, maintainable, and broadly compatible.