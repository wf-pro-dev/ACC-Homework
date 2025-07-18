# Eigenvalues, Eigenvectors, and Diagonalization: A Core Concept in Linear Algebra

This session delves into **eigenvalues** and **eigenvectors**, foundational concepts crucial for understanding linear transformations and their broad applications in science, engineering, and economics.

## What are Eigenvectors and Eigenvalues?

A **linear transformation** takes a vector as input and outputs another vector, often changing both its magnitude and direction. However, certain special vectors exist:

*   **Eigenvectors** (*x*): These are non-zero vectors that, when transformed by a matrix *A*, only change their *magnitude* (stretching or shrinking), not their *direction*. They stay along the same line.
*   **Eigenvalues** (*λ*): This is the scalar factor by which an eigenvector is scaled. It quantifies the change in magnitude.

Mathematically, this relationship is defined by the equation:
```
Ax = λx
```
where *A* is an *n* x *n* matrix, *x* is a non-zero eigenvector, and *λ* is its corresponding eigenvalue.

## Finding Eigenvalues: The Characteristic Equation

To find these special values and vectors, we manipulate the defining equation:

1.  Rearrange: `Ax - λx = 0`
2.  Introduce the Identity Matrix (*I*): `Ax - λIx = 0`
3.  Factor out *x*: `(A - λI)x = 0`

For a non-zero vector *x* to satisfy this equation, the matrix `(A - λI)` must be **singular** (meaning it does not have an inverse). A singular matrix has a determinant of zero.

This leads to the **Characteristic Equation**, the key to finding eigenvalues:

```
det(A - λI) = 0
```
Solving this polynomial equation for *λ* yields the eigenvalues.

## Example 1: Finding Eigenvalues and Eigenvectors for a 2x2 Matrix

Consider the matrix `A = [2 1; 1 2]`.

1.  **Form `(A - λI)`**:
    `A - λI = [ (2 - λ) 1; 1 (2 - λ) ]`

2.  **Calculate the determinant and set to zero**:
    `det(A - λI) = (2 - λ)(2 - λ) - (1)(1) = (2 - λ)² - 1 = 0`
    Expanding gives the **characteristic polynomial**: `λ² - 4λ + 3 = 0`

3.  **Solve for eigenvalues**:
    Factoring: `(λ - 1)(λ - 3) = 0`
    So, `λ₁ = 1` and `λ₂ = 3` are the eigenvalues.

4.  **Find corresponding eigenvectors** by solving `(A - λI)x = 0` for each eigenvalue:
    *   **For `λ₁ = 1`**: Solve `(A - 1I)x = 0`
        `[1 1; 1 1] [x₁; x₂] = [0; 0]`
        This implies `x₁ + x₂ = 0`, or `x₁ = -x₂`. Choosing `x₂ = 1`, the eigenvector `v₁ = [-1; 1]`.
    *   **For `λ₂ = 3`**: Solve `(A - 3I)x = 0`
        `[-1 1; 1 -1] [x₁; x₂] = [0; 0]`
        This implies `-x₁ + x₂ = 0`, or `x₁ = x₂`. Choosing `x₂ = 1`, the eigenvector `v₂ = [1; 1]`.

**Geometric Interpretation**: Vectors along the `[-1; 1]` direction are scaled by 1 (i.e., unchanged). Vectors along the `[1; 1]` direction are scaled by 3.

## Example 2: Handling Multiplicity in a 3x3 Matrix

Consider `B = [3 0 0; 0 1 -2; 0 -2 1]`.

1.  **Calculate `det(B - λI) = 0`**:
    `(3 - λ) * det([ (1 - λ) -2; -2 (1 - λ) ]) = 0`
    `(3 - λ) * ((1 - λ)² - 4) = 0`

2.  **Solve for eigenvalues**:
    From `(3 - λ) = 0`, we get `λ = 3`.
    From `(1 - λ)² - 4 = 0`, we get `(1 - λ)² = 4`, leading to `1 - λ = ±2`. This gives `λ = -1` and `λ = 3`.
    Thus, eigenvalues are `λ₁ = 3` (with **algebraic multiplicity** 2) and `λ₂ = -1` (with algebraic multiplicity 1).

3.  **Find eigenvectors**:
    *   For `λ = -1`: `v₁ = [0; 1; 1]`
    *   For `λ = 3`: `v₂ = [1; 0; 0]` and `v₃ = [0; -1; 1]` (two linearly independent eigenvectors, matching its **geometric multiplicity**).

## Summary of Steps to Find Eigenvalues and Eigenvectors

1.  **Form** the matrix `(A - λI)`. 
2.  **Calculate** the determinant of `(A - λI)` and set it to zero (`det(A - λI) = 0`) to obtain the **characteristic equation**.
3.  **Solve** the characteristic equation for `λ`. These are your **eigenvalues**.
4.  For each eigenvalue `λ`, **solve** the system `(A - λI)x = 0` for non-zero `x`. The solution space for `x` is the **eigenspace**, and any non-zero vector in it is an **eigenvector**.

## Properties of Eigenvalues

Two useful properties for checking calculations:

*   The **sum of the eigenvalues** equals the **trace** of the matrix (sum of main diagonal elements).
*   The **product of the eigenvalues** equals the **determinant** of the matrix.

## Matrix Diagonalization

A matrix `A` is **diagonalizable** if it is similar to a diagonal matrix `D`. This means there exists an invertible matrix *P* such that:

```
A = P D P⁻¹
```

*   `D` is a diagonal matrix containing the **eigenvalues** of `A` on its main diagonal.
*   `P` (the *modal matrix*) has the corresponding **eigenvectors** of `A` as its columns.

**Why is it useful?** Diagonal matrices are easy to work with. For instance, computing powers of `A` simplifies dramatically: `A^k = P D^k P⁻¹`, where `D^k` simply involves raising each diagonal entry of `D` to the power of `k`.

**Condition for Diagonalization**: An *n* x *n* matrix `A` is diagonalizable **if and only if** it has *n* **linearly independent eigenvectors**.

*   **Defective Matrices**: If, for any eigenvalue, its geometric multiplicity is less than its algebraic multiplicity (meaning you can't find enough linearly independent eigenvectors), the matrix is *not* diagonalizable. Example: `C = [1 1; 0 1]`.
*   **Special Cases**: 
    *   **Symmetric matrices** are *always* diagonalizable (and orthogonally diagonalizable).
    *   Matrices with **distinct eigenvalues** are *always* diagonalizable.

## Applications of Eigenvalues and Eigenvectors

These concepts are not just theoretical; they have profound applications:

*   **Systems of Linear Differential Equations**: Simplifying coupled systems into decoupled ones.
*   **Markov Chains**: Analyzing long-term behavior and steady-state distributions of systems.
*   **Principal Component Analysis (PCA)** in Data Science: Identifying main directions of variance in data for dimensionality reduction and feature extraction.
*   **Physics and Mechanics**: Describing principal axes of inertia, normal modes of vibration, and quantum mechanical states.

## Conclusion

Eigenvalues, eigenvectors, and diagonalization provide a powerful framework for understanding and simplifying linear transformations. By identifying special directions (eigenvectors) along which transformations merely scale (by eigenvalues), we gain deep insights into complex systems. Mastery of these concepts is fundamental for advanced topics in linear algebra and its widespread practical applications.